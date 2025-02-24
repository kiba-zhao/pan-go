package remotenode

import (
	"bytes"
	"context"
	"errors"
	appnode "pan/app/app_node"
	nodeitem "pan/extfs/node_item"
	remoteitem "pan/extfs/remote_item"
	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var ErrRemoteNodeFUSENameConflict = errors.New("remotefile.FUSERemoteNode Error: Name Conflict")

type RemoteNodeServiceFUSEProvider interface {
	LocalName() string
	RemoteNodeService() *RemoteNodeService
	remoteitem.RemoteItemServiceFUSEProvider
	nodeitem.NodeItemServiceFUSEProvider
}

type FUSERemoteNode struct {
	fs.Inode
	provider RemoteNodeServiceFUSEProvider
}

func (fusern *FUSERemoteNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

func (fusern *FUSERemoteNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {

	dirs := make([]fuse.DirEntry, 0)
	localEntry := fuse.DirEntry{Name: fusern.provider.LocalName(), Mode: fuse.S_IFDIR}
	dirs = append(dirs, localEntry)

	remoteNodeService := fusern.provider.RemoteNodeService()
	_, remotes, err := remoteNodeService.SelectAll()
	if err != nil {
		return nil, syscall.ENOENT
	}

	names := make([]string, 0)
	names = append(names, localEntry.Name)
	for _, remote := range remotes {
		idx, ok := slices.BinarySearch(names, remote.Name)
		if ok {
			err = ErrRemoteNodeFUSENameConflict
			break
		}
		names = slices.Insert(names, idx, remote.Name)
		dirEntry := fuse.DirEntry{Name: remote.Name, Mode: fuse.S_IFDIR}
		dirs = append(dirs, dirEntry)
	}

	releaseNames := make([]string, 0)
	for n := range fusern.Children() {
		if _, ok := slices.BinarySearch(names, n); ok {
			continue
		}
		releaseNames = append(releaseNames, n)
	}
	if len(releaseNames) > 0 {
		fusern.RmChild(releaseNames...)
	}
	if err != nil {
		return nil, syscall.ENOENT
	}

	return fs.NewListDirStream(dirs), 0
}

func (fusern *FUSERemoteNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	inode := fusern.GetChild(name)
	if name == fusern.provider.LocalName() {
		if inode != nil {
			_, ok := inode.Operations().(*nodeitem.FUSENodeItem)
			if ok {
				return inode, fs.OK
			}
			fusern.RmChild(name)
		}

		node := nodeitem.NewFUSENode(fusern.provider)
		inode = fusern.NewInode(ctx, node, fs.StableAttr{Mode: fuse.S_IFDIR})
		return inode, fs.OK
	}

	remoteNodeService := fusern.provider.RemoteNodeService()
	remote, err := remoteNodeService.SelectByName(name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	peerId, err := appnode.DecodePeerID(remote.PeerID)
	if err != nil {
		return nil, syscall.ENOENT
	}

	if inode != nil {
		fuseRemoteItemList, ok := inode.Operations().(*remoteitem.FUSERemoteItemList)
		if ok && bytes.Equal(fuseRemoteItemList.PeerID(), peerId) {
			return inode, fs.OK
		}
		fusern.RmChild(name)
	}

	node := remoteitem.NewFUSENode(peerId, fusern.provider)
	inode = fusern.NewInode(ctx, node, fs.StableAttr{Mode: fuse.S_IFDIR})
	return inode, fs.OK
}

func NewFUSENode(provider RemoteNodeServiceFUSEProvider) fs.InodeEmbedder {
	var node FUSERemoteNode
	node.provider = provider
	return &node
}
