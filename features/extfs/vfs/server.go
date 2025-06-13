package vfs

import (
	"context"
	"pan/lib/log"
	"sync"
)

// VFSFS is a virtual file system
type VFSFS interface {
	// Mount mounts the virtual file system using the given settings.
	// It returns an error if the settings are invalid or if the mount fails.
	Mount() error
	// Unmount unmounts the virtual file system.
	// It returns an error if the unmount fails.
	Unmount() error
}

type stdVFSServer struct {
	logger  log.Logger
	runtime *stdVFSFSRuntime

	mountPath   string
	mountPathRW sync.RWMutex

	hostName   string
	hostNameRW sync.RWMutex

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool
}

func (server *stdVFSServer) HostName() string {
	server.hostNameRW.RLock()
	defer server.hostNameRW.RUnlock()
	return server.hostName
}

func (server *stdVFSServer) SetHostName(hostName string) {
	server.hostNameRW.Lock()
	defer server.hostNameRW.Unlock()
	if server.hostName == hostName {
		return
	}
	server.hostName = hostName
}

func (server *stdVFSServer) MountPath() string {
	server.mountPathRW.RLock()
	defer server.mountPathRW.RUnlock()
	return server.mountPath
}

func (server *stdVFSServer) SetMountPath(mountPath string) {
	server.mountPathRW.Lock()
	defer server.mountPathRW.Unlock()
	if server.mountPath == mountPath {
		return
	}
	server.mountPath = mountPath
	server.Reload()
}

func (server *stdVFSServer) Reload() {
	server.logger.Debug("VFSServer", "Reload")

	server.reloadLock.Lock()
	defer server.reloadLock.Unlock()
	if server.reload {
		return
	}

	server.reload = true
	server.reloadChan <- struct{}{}
}

func (server *stdVFSServer) RunAndServe(ctx context.Context) error {
	server.logger.Debug("VFSServer", "RunAndServe begin")
	defer server.logger.Debug("VFSServer", "RunAndServe end")

	var err error
	var closed bool
	var vfsfs VFSFS
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-server.reloadChan:
			server.reloadLock.Lock()
			server.reload = false
			server.reloadLock.Unlock()
		}

		if vfsfs != nil {
			vfsfs.Unmount()
			vfsfs = nil
		}

		if closed {
			break
		}

		mountPath := server.MountPath()
		if len(mountPath) <= 0 {
			continue
		}
		vfsfs = NewVFSFS(server)
		err := vfsfs.Mount()
		if err != nil {
			vfsfs = nil
			server.logger.Error("VFSServer", "RunAndServe Error: "+err.Error())
		}
	}
	return err
}
