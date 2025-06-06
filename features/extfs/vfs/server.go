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

	settings   *VFSSettings
	settingsRW sync.RWMutex

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool
}

func (server *stdVFSServer) Settings() VFSSettings {
	server.settingsRW.RLock()
	defer server.settingsRW.RUnlock()
	return *server.settings
}

func (server *stdVFSServer) SetSettings(settings VFSSettings) {
	server.settingsRW.Lock()
	defer server.settingsRW.Unlock()
	*server.settings = settings
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

		settings := server.Settings()
		if !settings.Enabled {
			continue
		}
		vfsfs = NewVFSFS(server.runtime, settings)
		err := vfsfs.Mount()
		if err != nil {
			server.logger.Error("VFSServer", "RunAndServe Error: "+err.Error())
		}
	}
	return err
}
