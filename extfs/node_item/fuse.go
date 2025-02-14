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

type NodeItemServiceFUSEProvider interface {
	NodeItemService() *NodeItemService
}

type FUSENodeItem struct {
	fs.Inode
	provider NodeItemServiceFUSEProvider
}

func (fuseni *FUSENodeItem) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Mode = fuse.S_IFDIR
	out.Size = 4096
	now := time.Now()
	out.SetTimes(nil, &now, nil)
	out.Nlink = 1
	return fs.OK
}

func (fuseni *FUSENodeItem) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	dirs := make([]fuse.DirEntry, 0)
	names := make([]string, 0)
	nodeItemService := fuseni.provider.NodeItemService()
	err := nodeItemService.TraverseAll(func(nodeItem NodeItem) error {
		idx, ok := slices.BinarySearch(names, nodeItem.Name)
		if ok {
			return errors.New("nodeitem.FUSENodeItem Error: Name Conflict")
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

func (fuseni *FUSENodeItem) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	nodeItemService := fuseni.provider.NodeItemService()
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

func NewFUSENode(provider NodeItemServiceFUSEProvider) fs.InodeEmbedder {
	var fuse FUSENodeItem
	fuse.provider = provider
	return &fuse
}
