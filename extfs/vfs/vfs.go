package vfs

import (
	"context"
	appConfig "pan/app/config"
	"pan/app/injection"
	"pan/logger"
	"sync"
)

type VFSFileSystem interface {
	Mount(VFSSettings) error
	Unmount() error
}

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
	VFSFileSystem VFSFileSystem
	provider      injection.ComponentStoreProvider
	locker        sync.Mutex
	vfsSettings   *VFSSettings
	needReload    bool
	reloadChan    chan struct{}
	reloadOnce    sync.Once
	config        appConfig.Config[*VFSSettings]
}

func (vfs *VFS) ComponentStore() injection.ComponentStore {
	return vfs.provider.ComponentStore()
}

func (vfs *VFS) Components() []injection.Component {
	components := []injection.Component{
		injection.NewComponent(vfs, injection.ComponentNoneScope),
	}

	// append FUSEFileSystem
	var fuse FUSEFileSystem
	components = append(components,
		injection.NewComponent(&fuse, injection.ComponentNoneScope),
		injection.NewComponent[VFSFileSystem](&fuse, injection.ComponentInternalScope),
	)

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

		vfs.VFSFileSystem.Unmount()
		if closed {
			break
		}
		if !settings.Enabled {
			continue
		}

		err := vfs.VFSFileSystem.Mount(settings)
		if err != nil {
			logger.Default().Log(context.Background(), logger.LevelError, "extfs.vfs.VFS Error: %s", err.Error())
		}
	}
	return err
}
