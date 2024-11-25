package remotenode

import (
	"context"
	"encoding/base64"
	"errors"
	remoteitem "pan/extfs/remote_item"
	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

var ErrRemoteNodeFUSENameConflict = errors.New("remotefile.FUSERemoteNode Error: Name Conflict")

type RemoteNodeServiceFUSEProvider interface {
	RemoteNodeService() *RemoteNodeService
	remoteitem.RemoteItemServiceFUSEProvider
}

type FUSERemoteNode struct {
	fs.Inode
	Provider RemoteNodeServiceFUSEProvider
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
	remoteNodeService := fusern.Provider.RemoteNodeService()
	_, remotes, err := remoteNodeService.SelectAll()
	if err != nil {
		return nil, syscall.ENOENT
	}

	names := make([]string, 0)
	for _, remote := range remotes {
		idx, ok := slices.BinarySearch(names, remote.Name)
		if ok {
			err = ErrRemoteNodeFUSENameConflict
			break
		}
		names = slices.Insert(names, idx, remote.Name)
		dirs = append(dirs, fuse.DirEntry{Name: remote.Name, Mode: fuse.S_IFDIR})
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
	remoteNodeService := fusern.Provider.RemoteNodeService()
	remote, err := remoteNodeService.SelectByName(name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	inode := fusern.GetChild(name)
	if inode != nil {
		return inode, 0
	}

	peerId, err := base64.StdEncoding.DecodeString(remote.PeerID)
	if err != nil {
		return nil, syscall.ENOENT
	}

	remoteItem := &remoteitem.FUSERemoteItem{Provider: fusern.Provider, PeerID: peerId}
	inode = fusern.NewInode(ctx, remoteItem, fs.StableAttr{Mode: fuse.S_IFDIR})
	return inode, 0
}
