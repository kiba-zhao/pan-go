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
)

var ErrRemoteItemFUSENameConflict = errors.New("remoteitem.FUSERemoteItem Error: Name Conflict")

type FUSERemoteInfo interface {
	PeerID() peer.PeerID
}

type RemoteItemServiceFUSEProvider interface {
	RemoteItemService() *RemoteItemService
	RemoteFileInfoServiceFUSEProvider
}

type FUSERemoteItem struct {
	provider   RemoteItemServiceFUSEProvider
	remoteInfo FUSERemoteInfo
	itemId     uint
}

func (fuseri *FUSERemoteItem) PeerID() peer.PeerID {
	return fuseri.remoteInfo.PeerID()
}

func (fuseri *FUSERemoteItem) ItemID() uint {
	return fuseri.itemId
}

func (fuseri *FUSERemoteItem) GetRemoteFileAttr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	var condition RemoteItemRecordSelectCondition
	id := uint32(fuseri.ItemID())
	condition.ID = &id
	remoteItemService := fuseri.provider.RemoteItemService()
	remoteNodeItem, err := remoteItemService.SelectWithCondition(fuseri.PeerID(), &condition)
	if err == nil {

	}
	if err != nil {
		return syscall.ENOENT
	}

	mtime := time.Unix(remoteNodeItem.UpdatedAt, 0)
	out.SetTimes(nil, &mtime, nil)
	out.Size = uint64(remoteNodeItem.Size)
	out.Nlink = 1

	switch remoteNodeItem.FileType {
	case nodeitem.FileTypeFolder:
		out.Mode = fuse.S_IFDIR
	case nodeitem.FileTypeFile:
		out.Mode = fuse.S_IFREG
	}
	return 0
}

type FUSERemoteItemList struct {
	fs.Inode
	provider RemoteItemServiceFUSEProvider
	peerId   peer.PeerID
}

func (fusernil *FUSERemoteItemList) PeerID() peer.PeerID {
	return fusernil.peerId
}

func (fusernil *FUSERemoteItemList) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

func (fusernil *FUSERemoteItemList) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	names := make([]string, 0)
	dirs := make([]fuse.DirEntry, 0)
	remoteItemService := fusernil.provider.RemoteItemService()
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
	}, fusernil.PeerID())

	if err != nil {
		return nil, syscall.ENOENT
	}

	releaseNames := make([]string, 0)
	for n := range fusernil.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fusernil.RmChild(releaseNames...)
	}
	return fs.NewListDirStream(dirs), 0
}

func (fusernil *FUSERemoteItemList) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	var condition RemoteItemRecordSelectCondition
	condition.Name = &name
	remoteItemService := fusernil.provider.RemoteItemService()
	record, err := remoteItemService.SelectWithCondition(fusernil.PeerID(), &condition)
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

	inode := fusernil.GetChild(name)
	if inode != nil {
		fuseRemoteFile, ok := inode.Operations().(*FUSERemoteFile)
		if ok && fuseRemoteFile.ItemID() == uint(record.ID) && inode.Mode() == mode {
			return inode, fs.OK
		}
		fusernil.RmChild(name)
	}

	itemInfo := &FUSERemoteItem{remoteInfo: fusernil, itemId: uint(record.ID), provider: fusernil.provider}
	remoteFile := &FUSERemoteFile{provider: fusernil.provider, itemInfo: itemInfo}
	inode = fusernil.NewInode(ctx, remoteFile, fs.StableAttr{Mode: mode})

	return inode, fs.OK
}

func NewFUSENode(peerId peer.PeerID, provider RemoteItemServiceFUSEProvider) fs.InodeEmbedder {
	var fuse FUSERemoteItemList
	fuse.peerId = peerId
	fuse.provider = provider
	return &fuse
}
