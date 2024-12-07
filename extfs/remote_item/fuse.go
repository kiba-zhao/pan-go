package remoteitem

import (
	"context"
	"errors"
	"pan/app/peer"

	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	nodeitem "pan/extfs/node_item"
	remotefile "pan/extfs/remote_file"
)

var ErrRemoteItemFUSENameConflict = errors.New("remotefile.FUSERemoteItem Error: Name Conflict")

type RemoteItemServiceFUSEProvider interface {
	RemoteItemService() *RemoteItemService
	remotefile.RemoteFileServiceFUSEProvider
}

type FUSERemoteFileInfo struct {
	provider RemoteItemServiceFUSEProvider
	peerId   peer.PeerID
	itemId   int32
}

func (fuserfe *FUSERemoteFileInfo) PeerID() peer.PeerID {
	return fuserfe.peerId
}

func (fuserfe *FUSERemoteFileInfo) ItemID() int32 {
	return fuserfe.itemId
}

func (fuserfe *FUSERemoteFileInfo) GetRemoteFileAttr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	var condition RemoteItemRecordSelectCondition
	id := fuserfe.ItemID()
	condition.ID = &id
	remoteItemService := fuserfe.provider.RemoteItemService()
	remoteNode, err := remoteItemService.SelectWithCondition(fuserfe.peerId, &condition)
	if err == nil {
		mtime := time.Unix(remoteNode.UpdatedAt, 0)
		out.SetTimes(nil, &mtime, nil)
		out.Size = uint64(remoteNode.Size)
	}
	if err != nil {
		return syscall.ENOENT
	}

	out.Nlink = 1
	return 0
}

type FUSERemoteItem struct {
	fs.Inode
	Provider RemoteItemServiceFUSEProvider
	PeerID   peer.PeerID
}

func (fuserni *FUSERemoteItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

func (fuserni *FUSERemoteItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	names := make([]string, 0)
	dirs := make([]fuse.DirEntry, 0)
	remoteItemService := fuserni.Provider.RemoteItemService()
	err := remoteItemService.TraverseRecordWithPeerID(func(record *RemoteItemRecord) error {
		idx, ok := slices.BinarySearch(names, record.Name)
		if ok {
			return ErrRemoteItemFUSENameConflict
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
	}, fuserni.PeerID)

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

func (fuserni *FUSERemoteItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	var condition RemoteItemRecordSelectCondition
	condition.Name = &name
	remoteItemService := fuserni.Provider.RemoteItemService()
	record, err := remoteItemService.SelectWithCondition(fuserni.PeerID, &condition)
	if err != nil {
		return nil, syscall.ENOENT
	}

	var mode uint32
	switch record.FileType {
	case nodeitem.FileTypeFile:
		mode = fuse.S_IFREG
	case nodeitem.FileTypeFolder:
		mode = fuse.S_IFDIR
	}

	inode := fuserni.GetChild(name)
	if inode != nil {
		fuseRemoteFile, ok := inode.Operations().(*remotefile.FUSERemoteFile)
		if ok && fuseRemoteFile.ItemID() == record.ID && inode.Mode() == mode {
			return inode, fs.OK
		}
		fuserni.RmChild(name)
	}

	fileInfo := &FUSERemoteFileInfo{peerId: fuserni.PeerID, itemId: record.ID, provider: fuserni.Provider}
	remoteFile := &remotefile.FUSERemoteFile{Provider: fuserni.Provider, FUSERemoteFIleInfo: fileInfo}
	inode = fuserni.NewInode(ctx, remoteFile, fs.StableAttr{Mode: mode})

	return inode, fs.OK
}
