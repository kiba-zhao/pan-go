package vfs

import (
	"context"
	"errors"
	"os"

	"sync"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	nodeitem "pan/extfs/node_item"
	remoteblock "pan/extfs/remote_block"
	remotefile "pan/extfs/remote_file"
	remoteitem "pan/extfs/remote_item"
	remotenode "pan/extfs/remote_node"
)

var ErrFUSEUnavailable = errors.New("vfs.FUSE Error: Unavailable")

type FUSEServiceProvider struct {
	FUSEFileSystem *FUSEFileSystem
}

func (fusesp *FUSEServiceProvider) RemoteBlockService() *remoteblock.RemoteBlockService {
	return fusesp.FUSEFileSystem.RemoteBlockService
}

func (fusesp *FUSEServiceProvider) RemoteFileService() *remotefile.RemoteFileService {
	return fusesp.FUSEFileSystem.RemoteFileService
}

func (fusesp *FUSEServiceProvider) RemoteItemService() *remoteitem.RemoteItemService {
	return fusesp.FUSEFileSystem.RemoteItemService
}

func (fusesp *FUSEServiceProvider) RemoteNodeService() *remotenode.RemoteNodeService {
	return fusesp.FUSEFileSystem.RemoteNodeService
}

func (fusesp *FUSEServiceProvider) NodeItemService() *nodeitem.NodeItemService {
	return fusesp.FUSEFileSystem.NodeItemService
}

type FUSEFileSystem struct {
	RemoteBlockService *remoteblock.RemoteBlockService
	RemoteFileService  *remotefile.RemoteFileService
	RemoteItemService  *remoteitem.RemoteItemService
	RemoteNodeService  *remotenode.RemoteNodeService
	NodeItemService    *nodeitem.NodeItemService
	fs.Inode           `inject:"-"`
	server             *fuse.Server
	settings           *VFSSettings
	rw                 sync.RWMutex
	provider           *FUSEServiceProvider
	once               sync.Once
}

func (fusefs *FUSEFileSystem) Provider() *FUSEServiceProvider {
	fusefs.once.Do(func() {
		fusefs.provider = &FUSEServiceProvider{FUSEFileSystem: fusefs}
	})
	return fusefs.provider
}

func (fusefs *FUSEFileSystem) Mount(settings VFSSettings) error {
	fusefs.rw.Lock()
	defer fusefs.rw.Unlock()
	if fusefs.server != nil {
		return ErrFUSEUnavailable
	}

	opts := &fs.Options{}
	mountPath := settings.MountPath
	if err := os.MkdirAll(mountPath, 0755); err != nil {
		return err
	}
	fusefs.settings = &settings
	rawFS := fs.NewNodeFS(fusefs, opts)
	server, err := fuse.NewServer(rawFS, mountPath, &opts.MountOptions)
	if err != nil {
		return err
	}

	fusefs.server = server

	go server.Serve()
	err = server.WaitMount()
	return err
}

func (fusefs *FUSEFileSystem) Unmount() error {

	fusefs.rw.RLock()
	defer fusefs.rw.RUnlock()
	if fusefs.server == nil {
		return nil
	}

	err := fusefs.server.Unmount()
	if err == nil {
		fusefs.server.Wait()
		fusefs.server = nil
	}
	return err
}

func (fusefs *FUSEFileSystem) OnAdd(ctx context.Context) {
	var inode *fs.Inode
	if inode = fusefs.GetChild(fusefs.settings.LocalDirName); inode == nil {
		inode = fusefs.NewPersistentInode(ctx, &nodeitem.FUSENodeItem{Provider: fusefs.Provider()}, fs.StableAttr{Mode: syscall.S_IFDIR})
		fusefs.AddChild(fusefs.settings.LocalDirName, inode, true)
	}

	if inode = fusefs.GetChild(fusefs.settings.RemoteDirName); inode == nil {
		inode = fusefs.NewPersistentInode(ctx, &remotenode.FUSERemoteNode{Provider: fusefs.Provider()}, fs.StableAttr{Mode: syscall.S_IFDIR})
		fusefs.AddChild(fusefs.settings.RemoteDirName, inode, true)
	}

}
