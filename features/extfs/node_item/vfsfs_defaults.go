//go:build linux || (darwin && amd64)

package nodeitem

import (
	"context"
	"errors"
	"slices"
	"syscall"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type VFSFUSENodeItem struct {
	fs.Inode
	runtime VFSFUSENodeItemRuntime
}

// Getattr returns the attributes of the current node item.
// It always returns a directory with a size of 4096 bytes, and the
// current time as the last modified time. The number of hard links to
// the file is always 1.
// It implements the Getattr method of the fs.Inode interface.
func (fuseni *VFSFUSENodeItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

// Readdir returns a directory stream containing the node items as subdirectories.
// If any node item has the same name as an existing child node, an error is returned.
// The returned directory stream is released by calling RmChild on the VFSFUSENodeItem instance with
// the names of the node items that are not in the returned stream.
// It implements the Readdir method of the fs.Inode interface.
func (fuseni *VFSFUSENodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	dirs := make([]fuse.DirEntry, 0)
	names := make([]string, 0)
	nodeItemService := fuseni.runtime.NodeItemService()
	err := nodeItemService.TraverseAll(func(nodeItem NodeItem) error {
		idx, ok := slices.BinarySearch(names, nodeItem.Name)
		if ok {
			return errors.New("nodeitem.VFSFUSENodeItem Error: Name Conflict")
		}
		names = slices.Insert(names, idx, nodeItem.Name)

		dirEntry := fuse.DirEntry{}
		dirEntry.Name = nodeItem.Name
		switch nodeItem.FileType {
		case FileTypeFolder:
			dirEntry.Mode = fuse.S_IFDIR
		case FileTypeFile:
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

// Lookup searches for a child inode with the given name in the current node item.
// If the name matches the local node, it ensures the inode represents a local node
// and returns it. If the name corresponds to a remote node, it retrieves the
// corresponding remote node information and creates a new inode for it if necessary.
// Returns the found inode and fs.OK on success, or nil and ENOENT if the name
// does not correspond to any known node.
// It implements the Lookup method of the fs.Inode interface.
func (fuseni *VFSFUSENodeItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	nodeItemService := fuseni.runtime.NodeItemService()
	nodeItem, err := nodeItemService.SelectByName(name)
	if err != nil || !nodeItem.Available {
		return nil, syscall.ENOENT
	}

	var mode uint32
	switch nodeItem.FileType {
	case FileTypeFile:
		mode = fuse.S_IFREG
	case FileTypeFolder:
		mode = fuse.S_IFDIR
	}

	inode := fuseni.GetChild(name)
	if inode != nil {
		loopbackNode, ok := inode.Operations().(*fs.LoopbackNode)
		if ok && loopbackNode.RootData.Path == nodeItem.FilePath && inode.Mode() == mode {
			return inode, fs.OK
		}
		fuseni.RmChild(name)
	}

	itemNode := &fs.LoopbackNode{
		RootData: &fs.LoopbackRoot{
			Path: nodeItem.FilePath,
		},
	}
	itemNode.RootData.RootNode = itemNode
	inode = fuseni.NewInode(ctx, itemNode, fs.StableAttr{Mode: mode})

	return inode, fs.OK
}

// NewFUSENode creates a new VFSFUSENodeItem with the provided VFSFUSENodeItemRuntime.
// It initializes the VFSFUSENodeItem's runtime field with the given runtime
// and returns a pointer to the new VFSFUSENodeItem as an fs.InodeEmbedder.

func NewVFSFUSENodeItem(runtime VFSFUSENodeItemRuntime) fs.InodeEmbedder {
	var fuse VFSFUSENodeItem
	fuse.runtime = runtime
	return &fuse
}
