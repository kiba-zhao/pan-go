package vfs

import (
	"errors"
	"os"

	"sync"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"

	nodeitem "pan/extfs/node_item"
	remoteitem "pan/extfs/remote_item"
	remotenode "pan/extfs/remote_node"
)

var ErrFUSEUnavailable = errors.New("vfs.FUSE Error: Unavailable")

type FUSEServiceProvider struct {
	FUSEFileSystem *FUSEFileSystem
}

func (fusesp *FUSEServiceProvider) LocalName() string {
	return fusesp.FUSEFileSystem.settings.LocalName
}

func (fusesp *FUSEServiceProvider) RemoteFileStreamService() *remoteitem.RemoteFileStreamService {
	return fusesp.FUSEFileSystem.RemoteFileStreamService
}

func (fusesp *FUSEServiceProvider) RemoteFileInfoService() *remoteitem.RemoteFileInfoService {
	return fusesp.FUSEFileSystem.RemoteFileInfoService
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
	RemoteFileStreamService *remoteitem.RemoteFileStreamService
	RemoteFileInfoService   *remoteitem.RemoteFileInfoService
	RemoteItemService       *remoteitem.RemoteItemService
	RemoteNodeService       *remotenode.RemoteNodeService
	NodeItemService         *nodeitem.NodeItemService
	server                  *fuse.Server
	settings                *VFSSettings
	rw                      sync.RWMutex
	root                    fs.InodeEmbedder
	once                    sync.Once
}

func (fusefs *FUSEFileSystem) Root() fs.InodeEmbedder {
	fusefs.once.Do(func() {
		provider := &FUSEServiceProvider{FUSEFileSystem: fusefs}
		fusefs.root = remotenode.NewFUSENode(provider)
	})
	return fusefs.root
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
	rawFS := fs.NewNodeFS(fusefs.Root(), opts)
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
