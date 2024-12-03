package remotefile

import (
	"context"
	"errors"
	"pan/app/peer"
	remoteblock "pan/extfs/remote_block"
	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	nodeitem "pan/extfs/node_item"
)

var ErrRemoteFileFUSENameConflict = errors.New("remotefile.FUSERemoteFile Error: Name Conflict")

type FUSERemoteFIleInfo interface {
	PeerID() peer.PeerID
	ItemID() int32
	GetRemoteFileAttr(context.Context, *fuse.AttrOut) syscall.Errno
}

type RemoteFileServiceFUSEProvider interface {
	RemoteFileService() *RemoteFileService
	remoteblock.RemoteBlockServiceFUSEProvider
}

type FUSERemoteFile struct {
	fs.Inode
	Provider           RemoteFileServiceFUSEProvider
	Record             *RemoteFileRecord
	FUSERemoteFIleInfo FUSERemoteFIleInfo
}

func (fuserfe *FUSERemoteFile) PeerID() peer.PeerID {
	return fuserfe.FUSERemoteFIleInfo.PeerID()
}

func (fuserfe *FUSERemoteFile) ItemID() int32 {
	return fuserfe.FUSERemoteFIleInfo.ItemID()
}

func (fuserfe *FUSERemoteFile) ParentPath() string {
	if fuserfe.Record != nil {
		if fuserfe.Record.FileType == nodeitem.FileTypeFile {
			return fuserfe.Record.ParentPath
		}
		if fuserfe.Record.FileType == nodeitem.FileTypeFolder {
			return fuserfe.Record.FilePath
		}
	}
	return ""
}

func (fuserfe *FUSERemoteFile) Name() string {
	if fuserfe.Record != nil && fuserfe.Record.FileType == nodeitem.FileTypeFile {
		return fuserfe.Record.Name
	}
	return ""
}

func (fuserfe *FUSERemoteFile) GetRemoteFileAttr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	var condition RemoteFileRecordSelectCondition
	condition.ItemID = fuserfe.ItemID()
	condition.ParentPath = fuserfe.Record.ParentPath
	condition.Name = fuserfe.Record.Name
	remoteFileService := fuserfe.Provider.RemoteFileService()
	fileItem, err := remoteFileService.SelectWithCondition(fuserfe.PeerID(), &condition)
	if err == nil {
		mtime := time.Unix(fileItem.UpdatedAt, 0)
		out.SetTimes(nil, &mtime, nil)
		out.Size = uint64(fileItem.Size)
	}

	if err != nil {
		return syscall.ENOENT
	}

	out.Nlink = 1
	return 0
}

func (fuserfe *FUSERemoteFile) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if fuserfe.Record != nil {
		return fuserfe.GetRemoteFileAttr(ctx, out)
	}
	return fuserfe.FUSERemoteFIleInfo.GetRemoteFileAttr(ctx, out)
}

func (fuserfe *FUSERemoteFile) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {

	if fuserfe.Record.FileType != nodeitem.FileTypeFile {
		return nil, 0, syscall.ENFILE
	}

	// TODO: support writing
	if flags&(syscall.O_RDWR|syscall.O_WRONLY) != 0 {
		return nil, 0, syscall.EROFS
	}
	//
	return &remoteblock.FUSERemoteBlockReader{Provider: fuserfe.Provider, FUSERemoteFile: fuserfe}, fuse.FOPEN_DIRECT_IO, 0

}

func (fuserfe *FUSERemoteFile) Statfs(ctx context.Context, out *fuse.StatfsOut) syscall.Errno {
	// s := syscall.Statfs_t{}
	// err := syscall.Statfs(n.path(), &s)
	// if err != nil {
	// 	return ToErrno(err)
	// }
	// out.FromStatfsT(&s)
	return fs.OK
}

func (fuserfe *FUSERemoteFile) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	if fuserfe.Record != nil && fuserfe.Record.FileType != nodeitem.FileTypeFolder {
		return nil, syscall.ENOTDIR
	}

	dirs := make([]fuse.DirEntry, 0)
	var condition RemoteFileRecordSearchCondition
	condition.ItemID = fuserfe.ItemID()
	parentPath := fuserfe.ParentPath()
	condition.ParentPath = &parentPath

	remoteFileService := fuserfe.Provider.RemoteFileService()
	names := make([]string, 0)
	err := remoteFileService.TraverseRecordWithPeerID(func(record *RemoteFileRecord) error {
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

	}, fuserfe.FUSERemoteFIleInfo.PeerID(), &condition)

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

func (fuserfe *FUSERemoteFile) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if fuserfe.Record != nil && fuserfe.Record.FileType != nodeitem.FileTypeFolder {
		return nil, syscall.ENOTDIR
	}

	var condition RemoteFileRecordSelectCondition
	condition.ItemID = fuserfe.ItemID()
	condition.ParentPath = fuserfe.ParentPath()
	condition.Name = name

	remoteFileService := fuserfe.Provider.RemoteFileService()
	record, err := remoteFileService.SelectWithCondition(fuserfe.FUSERemoteFIleInfo.PeerID(), &condition)
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

	if inode != nil && inode.Mode() != mode {
		fuserfe.RmChild(name)
		inode = nil
	}
	if inode == nil {
		inode = fuserfe.NewInode(ctx, &FUSERemoteFile{Provider: fuserfe.Provider, Record: record, FUSERemoteFIleInfo: fuserfe}, fs.StableAttr{Mode: mode})
	}
	return inode, fs.OK
}
