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
	"path"
	"slices"
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
		names := make([]string, 0)
		for _, remote := range remotes {
			idx, ok := slices.BinarySearch(names, remote.Name)
			if ok {
				err = constant.ErrInternalError
				break
			}
			names = slices.Insert(names, idx, remote.Name)
			dirs = append(dirs, fuse.DirEntry{Name: remote.Name, Mode: fuse.S_IFDIR})
		}

		releaseNames := make([]string, 0)
		for n := range fusefs.Children() {
			if n == fusefs.localName {
				continue
			}
			if _, ok := slices.BinarySearch(names, n); ok {
				continue
			}
			releaseNames = append(releaseNames, n)
		}
		if len(releaseNames) > 0 {
			fusefs.RmChild(releaseNames...)
		}
	}

	if err != nil {
		return nil, syscall.ENOENT
	}

	return fs.NewListDirStream(dirs), 0
}

func (fusefs *FUSEFileSystem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if name == fusefs.localName {
		inode := fusefs.GetChild(name)
		if inode == nil {
			inode = fusefs.NewInode(ctx, &FUSENodeItem{FileSystem: fusefs}, fs.StableAttr{Mode: fuse.S_IFDIR})
		}
		return inode, 0
	}

	remote, err := fusefs.RemoteNodeService.SelectByName(name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	inode := fusefs.GetChild(name)
	if inode != nil {
		return inode, 0
	}

	nodeId, err := base64.StdEncoding.DecodeString(remote.NodeID)
	if err != nil {
		return nil, syscall.ENOENT
	}

	remoteInode := fusefs.NewInode(ctx, &FUSERemoteNodeItem{FileSystem: fusefs, NodeID: nodeId}, fs.StableAttr{Mode: fuse.S_IFDIR})
	return remoteInode, 0

}

type FUSENodeItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
}

func (fuseni *FUSENodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	dirs := make([]fuse.DirEntry, 0)
	names := make([]string, 0)
	err := fuseni.FileSystem.NodeItemService.TraverseAll(func(nodeItem models.NodeItem) error {
		idx, ok := slices.BinarySearch(names, nodeItem.Name)
		if ok {
			return constant.ErrInternalError
		}
		names = slices.Insert(names, idx, nodeItem.Name)

		dirEntry := fuse.DirEntry{}
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

	releaseNames := make([]string, 0)
	for n := range fuseni.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fuseni.RmChild(releaseNames...)
	}
	return fs.NewListDirStream(dirs), 0
}

func (fuseni *FUSENodeItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	nodeItem, err := fuseni.FileSystem.NodeItemService.SelectByName(name)
	if err != nil || !nodeItem.Available {
		return nil, syscall.ENOENT
	}

	inode := fuseni.GetChild(name)
	if inode != nil {
		if nodeItem.FileType == services.FileTypeFolder && inode.Mode() == fuse.S_IFDIR {
			return inode, 0
		}
		if nodeItem.FileType == services.FileTypeFile && inode.Mode() == fuse.S_IFREG {
			return inode, 0
		}
		inode = nil
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

	inode = fuseni.NewInode(ctx, itemNode, fs.StableAttr{Mode: mode})

	return inode, 0
}

type FUSERemoteNodeItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
}

func (fuserni *FUSERemoteNodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	names := make([]string, 0)
	dirs := make([]fuse.DirEntry, 0)
	err := fuserni.FileSystem.RemoteNodeItemService.TraverseRecordWithNodeID(func(record *models.RemoteNodeItemRecord) error {
		idx, ok := slices.BinarySearch(names, record.Name)
		if ok {
			return constant.ErrInternalError
		}
		names = slices.Insert(names, idx, record.Name)

		dirEntry := fuse.DirEntry{}
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

	releaseNames := make([]string, 0)
	for n := range fuserni.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fuserni.RmChild(releaseNames...)
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

	inode := fuserni.GetChild(name)
	if record.FileType == services.FileTypeFile {
		if inode != nil && inode.Mode() != fuse.S_IFREG {
			fuserni.RmChild(name)
			inode = nil
		}
		if inode == nil {
			inode = fuserni.NewInode(ctx, &FUSERemoteFileItem{FileSystem: fuserni.FileSystem, NodeID: fuserni.NodeID, ItemID: record.ID}, fs.StableAttr{Mode: fuse.S_IFREG})
		}
	}
	if record.FileType == services.FileTypeFolder {
		if inode != nil && inode.Mode() != fuse.S_IFDIR {
			fuserni.RmChild(name)
			inode = nil
		}
		if inode == nil {
			inode = fuserni.NewInode(ctx, &FUSERemoteFolderItem{FileSystem: fuserni.FileSystem, NodeID: fuserni.NodeID, ItemID: record.ID}, fs.StableAttr{Mode: fuse.S_IFDIR})
		}
	}
	return inode, 0
}

type FUSERemoteFolderItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
	ItemID     int32
	ParentPath string
}

func (fuserfi *FUSERemoteFolderItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {

	var err error
	if len(fuserfi.ParentPath) > 0 {
		dir, filename := path.Split(fuserfi.ParentPath)
		var condition models.RemoteFileItemRecordSelectCondition
		condition.ItemID = fuserfi.ItemID
		condition.ParentPath = dir
		condition.Name = filename
		fileItem, err := fuserfi.FileSystem.RemoteFileItemService.SelectWithCondition(fuserfi.NodeID, &condition)
		if err == nil {
			mtime := time.Unix(fileItem.UpdatedAt, 0)
			out.SetTimes(nil, &mtime, nil)
			out.Size = uint64(fileItem.Size)
		}
	} else {
		var condition models.RemoteNodeItemRecordSelectCondition
		condition.ID = &fuserfi.ItemID
		nodeItem, err := fuserfi.FileSystem.RemoteNodeItemService.SelectWithCondition(fuserfi.NodeID, &condition)
		if err == nil {
			mtime := time.Unix(nodeItem.UpdatedAt, 0)
			out.SetTimes(nil, &mtime, nil)
			out.Size = uint64(nodeItem.Size)
		}
	}

	if err != nil {
		return syscall.ENOENT
	}

	out.Nlink = 1
	return 0
}

func (fuserfi *FUSERemoteFolderItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	var condition models.RemoteFileItemRecordSearchCondition
	condition.ItemID = fuserfi.ItemID
	condition.ParentPath = &fuserfi.ParentPath

	names := make([]string, 0)
	err := fuserfi.FileSystem.RemoteFileItemService.TraverseRecordWithNodeID(func(record *models.RemoteFileItemRecord) error {
		idx, ok := slices.BinarySearch(names, record.Name)
		if ok {
			return constant.ErrInternalError
		}
		names = slices.Insert(names, idx, record.Name)

		dirEntry := fuse.DirEntry{}
		dirEntry.Name = record.Name
		switch record.FileType {
		case services.FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case services.FileTypeFile:
			dirEntry.Mode = fuse.S_IFREG
		}

		dirs = append(dirs, dirEntry)
		return nil

	}, fuserfi.NodeID, &condition)

	if err != nil {
		return nil, syscall.ENOENT
	}

	releaseNames := make([]string, 0)
	for n := range fuserfi.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fuserfi.RmChild(releaseNames...)
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

	inode := fuserfi.GetChild(name)
	if record.FileType == services.FileTypeFile {
		if inode != nil && inode.Mode() != fuse.S_IFREG {
			fuserfi.RmChild(name)
			inode = nil
		}
		if inode == nil {
			inode = fuserfi.NewInode(ctx, &FUSERemoteFileItem{FileSystem: fuserfi.FileSystem, NodeID: fuserfi.NodeID, ItemID: record.ItemID, ParentPath: record.ParentPath, Name: record.Name}, fs.StableAttr{Mode: fuse.S_IFREG})
		}
	}
	if record.FileType == services.FileTypeFolder {
		if inode != nil && inode.Mode() != fuse.S_IFDIR {
			fuserfi.RmChild(name)
			inode = nil
		}
		if inode == nil {
			inode = fuserfi.NewInode(ctx, &FUSERemoteFolderItem{FileSystem: fuserfi.FileSystem, NodeID: fuserfi.NodeID, ItemID: record.ItemID, ParentPath: record.FilePath}, fs.StableAttr{Mode: fuse.S_IFDIR})
		}
	}

	return inode, 0
}

type FUSERemoteFileReader struct {
	fileItem *FUSERemoteFileItem
	locker   sync.Mutex
	reader   io.ReadCloser
	offset   int64
}

func (fuserfr *FUSERemoteFileReader) closeReader() error {
	reader := fuserfr.reader
	var err error
	if reader != nil {
		err = reader.Close()
		fuserfr.reader = nil
	}

	return err
}

func (fuserfr *FUSERemoteFileReader) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()
	if fuserfr.offset < 0 {
		return nil, syscall.ENOENT
	}

	if fuserfr.reader != nil && off != fuserfr.offset {
		fuserfr.closeReader()
	}

	if fuserfr.reader == nil {
		var condition models.RemoteFileBlockSelectCondition
		condition.ItemID = fuserfr.fileItem.ItemID
		condition.ParentPath = fuserfr.fileItem.ParentPath
		condition.Name = fuserfr.fileItem.Name
		condition.Offset = off

		reader, err := fuserfr.fileItem.FileSystem.RemoteFileBlockService.SelectWithCondition(fuserfr.fileItem.NodeID, &condition)
		if err != nil {
			return nil, syscall.ENOENT
		}
		fuserfr.reader = reader
	}

	n, err := fuserfr.reader.Read(dest)
	if err != nil || n == 0 {
		fuserfr.closeReader()
		return nil, syscall.ENOENT
	}

	limit := int64(len(dest))
	var buffer []byte
	if int64(n) >= limit {
		fuserfr.offset = off + limit
		buffer = dest
	} else {
		fuserfr.offset = off + int64(n)
		buffer = dest[:n]
	}
	return fuse.ReadResultData(buffer), fs.OK
}

func (fuserfr *FUSERemoteFileReader) Release(ctx context.Context) syscall.Errno {
	fuserfr.locker.Lock()
	defer fuserfr.locker.Unlock()
	fuserfr.offset = -1
	err := fuserfr.closeReader()
	if err != nil {
		return syscall.ENOENT
	}
	return 0
}

type FUSERemoteFileItem struct {
	fs.Inode
	FileSystem *FUSEFileSystem
	NodeID     appNode.NodeID
	ItemID     int32
	ParentPath string
	Name       string
}

func (fuserfe *FUSERemoteFileItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {

	var err error
	if len(fuserfe.ParentPath) > 0 || len(fuserfe.Name) > 0 {
		var condition models.RemoteFileItemRecordSelectCondition
		condition.ItemID = fuserfe.ItemID
		condition.ParentPath = fuserfe.ParentPath
		condition.Name = fuserfe.Name
		fileItem, err := fuserfe.FileSystem.RemoteFileItemService.SelectWithCondition(fuserfe.NodeID, &condition)
		if err == nil {
			mtime := time.Unix(fileItem.UpdatedAt, 0)
			out.SetTimes(nil, &mtime, nil)
			out.Size = uint64(fileItem.Size)
		}
	} else {
		var condition models.RemoteNodeItemRecordSelectCondition
		condition.ID = &fuserfe.ItemID
		nodeItem, err := fuserfe.FileSystem.RemoteNodeItemService.SelectWithCondition(fuserfe.NodeID, &condition)
		if err == nil {
			mtime := time.Unix(nodeItem.UpdatedAt, 0)
			out.SetTimes(nil, &mtime, nil)
			out.Size = uint64(nodeItem.Size)
		}
	}

	if err != nil {
		return syscall.ENOENT
	}

	out.Nlink = 1
	return 0
}

func (fuserfe *FUSERemoteFileItem) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	// TODO: support writing
	if flags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}
	//
	return &FUSERemoteFileReader{fileItem: fuserfe}, fuse.FOPEN_DIRECT_IO, 0

}
