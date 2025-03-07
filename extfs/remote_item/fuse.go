// Define fuse for remote item
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
	// Returns the peer ID of the remote item provider.
	PeerID() peer.PeerID
}

type RemoteItemServiceFUSEProvider interface {
	// Returns the remote item service.
	RemoteItemService() *RemoteItemService
	// It extends the RemoteFileInfoServiceFUSEProvider interface.
	RemoteFileInfoServiceFUSEProvider
}

type FUSERemoteItem struct {
	provider   RemoteItemServiceFUSEProvider
	remoteInfo FUSERemoteInfo
	itemId     uint
}

// PeerID returns the peer ID of the remote item provider.
func (fuseri *FUSERemoteItem) PeerID() peer.PeerID {
	return fuseri.remoteInfo.PeerID()
}

// ItemID returns the ID of the remote item.
func (fuseri *FUSERemoteItem) ItemID() uint {
	return fuseri.itemId
}

// GetRemoteFileAttr returns the file attributes of the remote item.
//
// It queries the remote item record by ID from the remote item service,
// and sets the attributes in the AttrOut structure.
//
// If the remote item record does not exist, it returns syscall.ENOENT.
// Otherwise, it returns 0.
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

// PeerID returns the peer ID of the remote item provider.
func (fusernil *FUSERemoteItemList) PeerID() peer.PeerID {
	return fusernil.peerId
}

// Getattr returns the attributes of the current remote item list.
//
// It always returns a directory with a size of 4096 bytes, and the
// current time as the last modified time. The number of hard links to
// the file is always 1.
//
// It implements the Getattr method of the fs.Inode interface.
func (fusernil *FUSERemoteItemList) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

// Readdir retrieves a directory stream containing entries of remote items
// associated with the current FUSERemoteItemList. It traverses the remote
// item records using the peer ID and creates a list of directory entries
// representing each item. If a name conflict occurs, an error is returned.
// The function also removes child nodes that are not in the retrieved list
// from the FUSERemoteItemList's children. Returns a newly created directory
// stream and syscall.ENOENT if an error occurs during traversal.

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

// Lookup searches for a child inode with the given name within the remote item list.
// It queries the remote item service for the record associated with the name
// and determines the file type to set the mode accordingly. If an inode with the
// specified name already exists and matches the expected file type and item ID,
// it returns the existing inode. Otherwise, it removes the existing child inode
// and creates a new one with the corresponding remote file information.
// Returns the found or newly created inode and fs.OK on success, or nil and ENOENT if
// the record with the specified name does not exist.

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

// NewFUSENode creates a new FUSENodeItem with the given peer ID and remote item service provider.
// It initializes the FUSENodeItem's peer ID and provider fields with the given peer ID and provider,
// and returns a pointer to the new FUSENodeItem as an fs.InodeEmbedder.
func NewFUSENode(peerId peer.PeerID, provider RemoteItemServiceFUSEProvider) fs.InodeEmbedder {
	var fuse FUSERemoteItemList
	fuse.peerId = peerId
	fuse.provider = provider
	return &fuse
}
