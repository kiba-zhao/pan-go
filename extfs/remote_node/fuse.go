// Define FUSE INode for Remote Node
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

// Remote Node FUSE Provider
type RemoteNodeServiceFUSEProvider interface {
	// LocalName returns the local name of the remote node
	LocalName() string
	// RemoteNodeService returns the remote node service.
	RemoteNodeService() *RemoteNodeService
	// It extends the following interfaces:
	remoteitem.RemoteItemServiceFUSEProvider
	nodeitem.NodeItemServiceFUSEProvider
}

type FUSERemoteNode struct {
	fs.Inode
	provider RemoteNodeServiceFUSEProvider
}

// Getattr returns the attributes of the root node of the remote node.
// It always returns a directory with a size of 4096 bytes, and the
// current time as the last modified time. The number of hard links to
// the file is always 1.
// It implements the Getattr method of the fs.Inode interface.
func (fusern *FUSERemoteNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

// Readdir returns a directory stream containing the local node and all
// remote nodes as subdirectories. If any remote node has the same name
// as the local node, an error is returned. The returned directory stream
// is released by calling RmChild on the FUSERemoteNode instance with
// the names of the remote nodes that are not in the returned stream.
// It implements the Readdir method of the fs.Inode interface.
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

// Lookup searches for a child inode with the given name in the current remote node.
// If the name matches the local node, it ensures the inode represents a local node
// and returns it. If the name corresponds to a remote node, it retrieves the
// corresponding remote node information and creates a new inode for it if necessary.
// Returns the found inode and fs.OK on success, or nil and ENOENT if the name
// does not correspond to any known node.
// It implements the Lookup method of the fs.Inode interface.
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

// NewFUSENode creates a new fs.InodeEmbedder for a remote node.
//
// provider is the RemoteNodeServiceFUSEProvider for this node.
func NewFUSENode(provider RemoteNodeServiceFUSEProvider) fs.InodeEmbedder {
	var node FUSERemoteNode
	node.provider = provider
	return &node
}
