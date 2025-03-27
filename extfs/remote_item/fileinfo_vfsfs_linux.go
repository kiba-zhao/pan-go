//go:build linux || (darwin && amd64)

package remoteitem

import (
	"context"
	"errors"
	nodeitem "pan/extfs/node_item"
	"path"
	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var ErrRemoteFileFUSENameConflict = errors.New("remoteitem.VFSFUSERemoteFileInfo Error: Name Conflict")

type VFSFUSERemoteItemInfo interface {
	VFSFUSERemoteInfo
	ItemID() uint
	GetRemoteFileAttr(context.Context, *fuse.AttrOut) syscall.Errno
}

type VFSFUSERemoteFile struct {
	fs.Inode
	runtime  VFSFUSERemoteFileRuntime
	itemInfo VFSFUSERemoteItemInfo
	filePath string
}

func (fuserfe *VFSFUSERemoteFile) ItemID() uint {
	return fuserfe.itemInfo.ItemID()
}

func (fuserfe *VFSFUSERemoteFile) GetRemoteFileAttr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {

	var condition RemoteFileInfoRecordSelectCondition
	condition.ItemID = uint32(fuserfe.itemInfo.ItemID())
	condition.FilePath = fuserfe.filePath
	remoteFileInfoService := fuserfe.runtime.RemoteFileInfoService()
	fileInfo, err := remoteFileInfoService.SelectWithCondition(fuserfe.itemInfo.PeerID(), &condition)

	if err != nil {
		return syscall.ENOENT
	}

	mtime := time.Unix(fileInfo.UpdatedAt, 0)
	out.SetTimes(nil, &mtime, nil)
	out.Size = uint64(fileInfo.Size)
	out.Nlink = 1

	switch fileInfo.FileType {
	case nodeitem.FileTypeFolder:
		out.Mode = fuse.S_IFDIR
	case nodeitem.FileTypeFile:
		out.Mode = fuse.S_IFREG
	}

	return 0
}

func (fuserfe *VFSFUSERemoteFile) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if len(fuserfe.filePath) > 0 {
		return fuserfe.GetRemoteFileAttr(ctx, out)
	}
	return fuserfe.itemInfo.GetRemoteFileAttr(ctx, out)
}

func (fuserfe *VFSFUSERemoteFile) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {

	// TODO: support writing
	if flags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}
	//

	remoteFileStreamService := fuserfe.runtime.RemoteFileStreamService()
	reader, err := remoteFileStreamService.ReadWithPeerID(fuserfe.itemInfo.PeerID(), fuserfe.itemInfo.ItemID(), fuserfe.filePath)
	if err != nil {
		return nil, 0, syscall.ENOENT
	}
	return &VFSFUSERemoteFileStreamReader{reader: reader}, fuse.FOPEN_DIRECT_IO, 0

}

func (fuserfe *VFSFUSERemoteFile) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	// s := syscall.Statfs_t{}
	// err := syscall.Statfs(n.path(), &s)
	// if err != nil {
	// 	return ToErrno(err)
	// }
	// out.FromStatfsT(&s)
	return fs.OK
}

func (fuserfe *VFSFUSERemoteFile) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	var condition RemoteFileInfoRecordSearchCondition
	condition.ItemID = uint32(fuserfe.itemInfo.ItemID())
	condition.ParentPath = fuserfe.filePath

	remoteFileInfoService := fuserfe.runtime.RemoteFileInfoService()
	names := make([]string, 0)
	err := remoteFileInfoService.TraverseRecordWithPeerID(func(record *RemoteFileInfoRecord) error {
		idx, ok := slices.BinarySearch(names, record.Name)
		if ok {
			return ErrRemoteFileFUSENameConflict
		}
		names = slices.Insert(names, idx, record.Name)

		dirEntry := fuse.DirEntry{}
		dirEntry.Name = record.Name
		switch record.FileType {
		case nodeitem.FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case nodeitem.FileTypeFile:
			dirEntry.Mode = fuse.S_IFREG
		}

		dirs = append(dirs, dirEntry)
		return nil

	}, fuserfe.itemInfo.PeerID(), &condition)

	if err != nil {
		return nil, syscall.ENOENT
	}

	releaseNames := make([]string, 0)
	for n := range fuserfe.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fuserfe.RmChild(releaseNames...)
	}
	return fs.NewListDirStream(dirs), 0
}

func (fuserfe *VFSFUSERemoteFile) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	var condition RemoteFileInfoRecordSelectCondition
	condition.ItemID = uint32(fuserfe.itemInfo.ItemID())
	if len(fuserfe.filePath) > 0 {
		condition.FilePath = path.Join(fuserfe.filePath, name)
	} else {
		condition.FilePath = name
	}

	remoteFileInfoService := fuserfe.runtime.RemoteFileInfoService()
	record, err := remoteFileInfoService.SelectWithCondition(fuserfe.itemInfo.PeerID(), &condition)
	if err != nil {
		return nil, syscall.ENOENT
	}

	inode := fuserfe.GetChild(name)
	var mode uint32
	switch record.FileType {
	case nodeitem.FileTypeFile:
		mode = fuse.S_IFREG
	case nodeitem.FileTypeFolder:
		mode = fuse.S_IFDIR
	}

	if inode != nil {
		if inode.Mode() == mode {
			return inode, fs.OK
		}
		fuserfe.RmChild(name)
	}

	inode = fuserfe.NewInode(ctx, &VFSFUSERemoteFile{runtime: fuserfe.runtime, itemInfo: fuserfe.itemInfo, filePath: condition.FilePath}, fs.StableAttr{Mode: mode})
	return inode, fs.OK
}
