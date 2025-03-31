package vfs

import (
	"context"
	appConfig "pan/lib/config"
	"pan/lib/injection"
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

// New creates a new virtual file system.
//
// It takes a component store provider to create the virtual file system.
// It will load the configuration from "extfs_vfs.toml" file.
// If the configuration file does not exist, it will panic with an error.
// If the configuration file exists but is invalid, it will panic with an error.
func New(provider injection.ComponentStoreProvider) interface{} {

	vfs := &VFS{}
	vfs.provider = provider

	config, err := appConfig.NewConfig[*VFSSettings]("extfs_vfs.toml")
	if err != nil {
		panic(err)
	}
	vfs.config = config

	return vfs
}

type VFS struct {
	provider     injection.ComponentStoreProvider
	locker       sync.Mutex
	vfsSettings  *VFSSettings
	needReload   bool
	reloadChan   chan struct{}
	reloadOnce   sync.Once
	config       appConfig.Config[*VFSSettings]
	vfsfsRuntime *VFSFSRuntime
	vfsfsmOnce   sync.Once
}

func (vfs *VFS) VFSFSRuntime() *VFSFSRuntime {
	vfs.vfsfsmOnce.Do(func() {
		vfs.vfsfsRuntime = &VFSFSRuntime{}
	})
	return vfs.vfsfsRuntime
}

func (vfs *VFS) ComponentStore() injection.ComponentStore {
	return vfs.provider.ComponentStore()
}

func (vfs *VFS) Components() []injection.Component {
	components := []injection.Component{
		injection.NewComponent(vfs, injection.ComponentNoneScope),
		injection.NewComponent(vfs.VFSFSRuntime(), injection.ComponentNoneScope),
	}

	return components
}

func (vfs *VFS) Modules() []interface{} {
	return []interface{}{
		vfs.config,
	}
}

func (vfs *VFS) OnConfigUpdated(settings appConfig.AppSettings) {
	vfs.locker.Lock()
	defer vfs.locker.Unlock()
	if vfs.config == nil {
		return
	}

	defaultsSettings := newDefaultsVFSSettings(settings)
	vfs.config.SetDefaults(defaultsSettings)
	vfsSettings, err := vfs.config.Load()
	if err != nil {
		panic(err)
	}
	vfs.vfsSettings = vfsSettings
	// trigger to reload
	if vfs.needReload {
		return
	}
	vfs.needReload = true
	vfs.ReloadChan() <- struct{}{}
}

func (vfs *VFS) ReloadChan() chan struct{} {
	vfs.reloadOnce.Do(func() {
		vfs.reloadChan = make(chan struct{}, 1)
	})
	return vfs.reloadChan
}

func (vfs *VFS) Ready(ctx context.Context) error {

	var vfsfs VFSFS
	var err error
	closed := false
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-vfs.ReloadChan():
		}

		vfs.locker.Lock()
		vfs.needReload = false
		settings := *vfs.vfsSettings
		vfs.locker.Unlock()

		if vfsfs != nil {
			vfsfs.Unmount()
			vfsfs = nil
		}

		if closed {
			break
		}

		if !settings.Enabled {
			continue
		}

		vfsfsRuntime := vfs.VFSFSRuntime()
		if vfsfsRuntime == nil {
			continue
		}

		vfsfs = NewVFSFS(vfsfsRuntime, settings)
		err := vfsfs.Mount()
		if err != nil {
			log.Default().Log(context.Background(), log.LevelError, "extfs.vfs.VFS Error: "+err.Error())
		}
		// if vfs.VFSFS != nil {
		// 	err := vfs.VFSFS.Mount(settings)
		// 	if err != nil {
		// 		logger.Default().Log(context.Background(), logger.LevelError, "extfs.vfs.VFS Error: %s", err.Error())
		// 	}
		// }
	}
	return err
}
