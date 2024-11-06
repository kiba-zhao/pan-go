package vfs

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"pan/app/constant"
	appNode "pan/app/node"
	"pan/extfs/models"
	"pan/extfs/services"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FUSEFileSystem struct {
	RemoteFileBlockService *services.RemoteFileBlockService
	RemoteFileItemService  *services.RemoteFileItemService
	RemoteNodeItemService  *services.RemoteNodeItemService
	RemoteNodeService      *services.RemoteNodeService
	NodeItemService        *services.NodeItemService
	fs.Inode               `inject:"-"`
	server                 *fuse.Server
	localName              string
	rw                     sync.RWMutex
}

func (fusefs *FUSEFileSystem) Mount(settings VFSSettings) error {
	fusefs.rw.Lock()
	defer fusefs.rw.Unlock()
	if fusefs.server != nil {
		return constant.ErrConflict
	}

	opts := &fs.Options{}
	mountPath := settings.MountPath
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return err
	}
	rawFS := fs.NewNodeFS(fusefs, opts)
	server, err := fuse.NewServer(rawFS, mountPath, &opts.MountOptions)
	if err != nil {
		return err
	}

	fusefs.server = server
	fusefs.localName = settings.LocalName

	go server.Serve()
	err = server.WaitMount()
	return err
}

func (fusefs *FUSEFileSystem) Unmount() error {
	fusefs.rw.RLock()
	defer fusefs.rw.RUnlock()
	if fusefs.server == nil {
		return nil
	}

	err := fusefs.server.Unmount()
	if err == nil {
		fusefs.server.Wait()
		fusefs.server = nil
	}
	return err
}

func (fusefs *FUSEFileSystem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)

	fusefs.rw.RLock()
	dirs = append(dirs, fuse.DirEntry{Name: fusefs.localName, Mode: fuse.S_IFDIR})
	fusefs.rw.RUnlock()

	_, remotes, err := fusefs.RemoteNodeService.SelectAll()
	if err == nil {
		for _, remote := range remotes {
			dirs = append(dirs, fuse.DirEntry{Name: remote.Name, Mode: fuse.S_IFDIR, Ino: uint64(remote.ID)})
		}
	}

	return fs.NewListDirStream(dirs), 0
}

func (fusefs *FUSEFileSystem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if name == fusefs.localName {
		localInode := fusefs.NewInode(ctx, &FUSENodeItem{FileSystem: fusefs}, fs.StableAttr{Mode: fuse.S_IFDIR})
		return localInode, 0
	}

	remote, err := fusefs.RemoteNodeService.SelectByName(name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	nodeId, err := base64.StdEncoding.DecodeString(remote.NodeID)
	if err != nil {
		return nil, syscall.ENOENT
	}

	remoteInode := fusefs.NewInode(ctx, &FUSERemoteNodeItem{FileSystem: fusefs, NodeID: nodeId}, fs.StableAttr{Ino: uint64(remote.ID), Mode: fuse.S_IFDIR})
	return remoteInode, 0

}

type FUSENodeItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
}

func (fuseni *FUSENodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	dirs := make([]fuse.DirEntry, 0)
	err := fuseni.FileSystem.NodeItemService.TraverseAll(func(nodeItem models.NodeItem) error {
		dirEntry := fuse.DirEntry{}
		dirEntry.Ino = uint64(nodeItem.ID)
		dirEntry.Name = nodeItem.Name
		switch nodeItem.FileType {
		case services.FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case services.FileTypeFile:
			dirEntry.Mode = fuse.S_IFREG
		}
		dirs = append(dirs, dirEntry)
		return nil
	})
	if err != nil {
		return nil, syscall.ENOENT
	}
	return fs.NewListDirStream(dirs), 0
}

func (fuseni *FUSENodeItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	nodeItem, err := fuseni.FileSystem.NodeItemService.SelectByName(name)
	if err != nil || !nodeItem.Available {
		return nil, syscall.ENOENT
	}

	itemNode := &fs.LoopbackNode{
		RootData: &fs.LoopbackRoot{
			Path: nodeItem.FilePath,
		},
	}
	itemNode.RootData.RootNode = itemNode

	var mode uint32
	switch nodeItem.FileType {
	case services.FileTypeFile:
		mode = fuse.S_IFREG
	case services.FileTypeFolder:
		mode = fuse.S_IFDIR
	}

	inode := fuseni.NewInode(ctx, itemNode, fs.StableAttr{Ino: uint64(nodeItem.ID), Mode: mode})

	return inode, 0
}

type FUSERemoteNodeItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
}

func (fuserni *FUSERemoteNodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	err := fuserni.FileSystem.RemoteNodeItemService.TraverseRecordWithNodeID(func(record *models.RemoteNodeItemRecord) error {
		dirEntry := fuse.DirEntry{}
		dirEntry.Ino = uint64(record.ID)
		dirEntry.Name = record.Name
		switch record.FileType {
		case services.FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case services.FileTypeFile:
			dirEntry.Mode = fuse.S_IFREG
		}
		dirs = append(dirs, dirEntry)
		return nil
	}, fuserni.NodeID)

	if err != nil {
		return nil, syscall.ENOENT
	}

	return fs.NewListDirStream(dirs), 0
}

func (fuserni *FUSERemoteNodeItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	var condition models.RemoteNodeItemRecordSelectCondition
	condition.Name = &name
	record, err := fuserni.FileSystem.RemoteNodeItemService.SelectWithCondition(fuserni.NodeID, &condition)
	if err != nil {
		return nil, syscall.ENOENT
	}

	mtime := time.Unix(record.UpdatedAt, 0)
	var remoteInode *fs.Inode
	if record.FileType == services.FileTypeFile {
		remoteInode = fuserni.NewInode(ctx, &FUSERemoteFileItem{FileSystem: fuserni.FileSystem, NodeID: fuserni.NodeID, ItemID: record.ID, Size: uint64(record.Size), MTime: mtime}, fs.StableAttr{Ino: uint64(record.ID), Mode: fuse.S_IFREG})
	}
	if record.FileType == services.FileTypeFolder {
		remoteInode = fuserni.NewInode(ctx, &FUSERemoteFolderItem{FileSystem: fuserni.FileSystem, NodeID: fuserni.NodeID, ItemID: record.ID, seq: 1, Size: uint64(record.Size), MTime: mtime}, fs.StableAttr{Ino: uint64(record.ID), Mode: fuse.S_IFDIR})
	}
	return remoteInode, 0
}

type FUSERemoteFileInfo struct {
	Name string
	Ino  uint64
}

type FUSERemoteFolderItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
	ItemID     int32
	ParentPath string
	Size       uint64
	MTime      time.Time
	seq        uint64
	fileInfos  []FUSERemoteFileInfo
	locker     sync.Mutex
}

func (fuserfi *FUSERemoteFolderItem) CompareFileInfo(fileInfo FUSERemoteFileInfo, target FUSERemoteFileInfo) int {
	return strings.Compare(target.Name, fileInfo.Name)
}

func (fuserfi *FUSERemoteFolderItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size = fuserfi.Size
	out.SetTimes(nil, &fuserfi.MTime, nil)
	out.Nlink = 1
	return 0
}

func (fuserfi *FUSERemoteFolderItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	var condition models.RemoteFileItemRecordSearchCondition
	condition.ItemID = fuserfi.ItemID
	condition.ParentPath = &fuserfi.ParentPath

	fuserfi.locker.Lock()
	fileInfos_ := make([]FUSERemoteFileInfo, 0)
	err := fuserfi.FileSystem.RemoteFileItemService.TraverseRecordWithNodeID(func(record *models.RemoteFileItemRecord) error {
		var fileInfo FUSERemoteFileInfo
		fileInfo.Name = record.Name

		if idx, ok := slices.BinarySearchFunc(fuserfi.fileInfos, fileInfo, fuserfi.CompareFileInfo); ok {
			fileInfo = fuserfi.fileInfos[idx]
		} else {
			fileInfo.Ino = fuserfi.seq
			fuserfi.seq++
		}

		idx, _ := slices.BinarySearchFunc(fileInfos_, fileInfo, fuserfi.CompareFileInfo)
		fileInfos_ = slices.Insert(fileInfos_, idx, fileInfo)

		dirEntry := fuse.DirEntry{}
		dirEntry.Name = record.Name
		dirEntry.Ino = fileInfo.Ino
		switch record.FileType {
		case services.FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case services.FileTypeFile:
			dirEntry.Mode = fuse.S_IFREG
		}

		dirs = append(dirs, dirEntry)
		return nil

	}, fuserfi.NodeID, &condition)
	fuserfi.fileInfos = fileInfos_
	fuserfi.locker.Unlock()

	if err != nil {
		return nil, syscall.ENOENT
	}

	return fs.NewListDirStream(dirs), 0
}

func (fuserfi *FUSERemoteFolderItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	var condition models.RemoteFileItemRecordSelectCondition
	condition.ItemID = fuserfi.ItemID
	condition.ParentPath = fuserfi.ParentPath
	condition.Name = name
	record, err := fuserfi.FileSystem.RemoteFileItemService.SelectWithCondition(fuserfi.NodeID, &condition)
	if err != nil {
		return nil, syscall.ENOENT
	}

	fuserfi.locker.Lock()
	var fileInfo FUSERemoteFileInfo
	fileInfo.Name = record.Name
	if idx, ok := slices.BinarySearchFunc(fuserfi.fileInfos, fileInfo, fuserfi.CompareFileInfo); !ok {
		fileInfo.Ino = fuserfi.seq
		fuserfi.seq++
		fuserfi.fileInfos = slices.Insert(fuserfi.fileInfos, idx, fileInfo)
	} else {
		fileInfo = fuserfi.fileInfos[idx]
	}
	fuserfi.locker.Unlock()

	mtime := time.Unix(record.UpdatedAt, 0)
	var remoteInode *fs.Inode
	if record.FileType == services.FileTypeFile {
		remoteInode = fuserfi.NewInode(ctx, &FUSERemoteFileItem{FileSystem: fuserfi.FileSystem, NodeID: fuserfi.NodeID, ItemID: record.ItemID, ParentPath: record.ParentPath, Name: record.Name, Size: uint64(record.Size), MTime: mtime}, fs.StableAttr{Ino: fileInfo.Ino, Mode: fuse.S_IFREG})
	}
	if record.FileType == services.FileTypeFolder {
		remoteInode = fuserfi.NewInode(ctx, &FUSERemoteFolderItem{FileSystem: fuserfi.FileSystem, NodeID: fuserfi.NodeID, ItemID: record.ItemID, ParentPath: record.FilePath, seq: 1, Size: uint64(record.Size), MTime: mtime}, fs.StableAttr{Ino: fileInfo.Ino, Mode: fuse.S_IFDIR})
	}

	return remoteInode, 0
}

type FUSERemoteFileItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
	ItemID     int32
	ParentPath string
	Name       string
	Size       uint64
	MTime      time.Time
	locker     sync.Mutex
	reader     io.Reader
	offset     int64
}

func (fuserfe *FUSERemoteFileItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size = fuserfe.Size
	out.SetTimes(nil, &fuserfe.MTime, nil)
	out.Nlink = 1
	return 0
}

func (fuserfe *FUSERemoteFileItem) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	// TODO: support writing
	if flags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}
	//
	return fuserfe, fuse.FOPEN_DIRECT_IO, 0
}

func (fuserfe *FUSERemoteFileItem) Read(ctx context.Context, fh fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fuserfe.locker.Lock()
	defer fuserfe.locker.Unlock()

	limit := int64(len(dest))
	if fuserfe.reader == nil || off != fuserfe.offset {
		var condition models.RemoteFileBlockSelectCondition
		condition.ItemID = fuserfe.ItemID
		condition.ParentPath = fuserfe.ParentPath
		condition.Name = fuserfe.Name
		condition.Offset = off
		condition.Limit = limit

		reader, err := fuserfe.FileSystem.RemoteFileBlockService.SelectWithCondition(fuserfe.NodeID, &condition)
		if err != nil {
			return nil, syscall.ENOENT
		}
		fuserfe.reader = reader
	}

	buffer, err := io.ReadAll(fuserfe.reader)
	if err != nil {
		return nil, syscall.ENOENT
	}

	end := int64(len(buffer))
	if end > limit {
		fuserfe.offset = off + limit
		return fuse.ReadResultData(buffer[:limit]), 0
	}
	fuserfe.offset = off + end
	return fuse.ReadResultData(buffer), 0
}
