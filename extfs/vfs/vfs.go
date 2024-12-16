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

func New(store injection.ComponentStore) interface{} {
	return &VFS{store: store}
}

type VFS struct {
	VFSFileSystem VFSFileSystem
	store         injection.ComponentStore
	locker        sync.Mutex
	vfsSettings   *VFSSettings
	needReload    bool
	reloadChan    chan struct{}
	reloadOnce    sync.Once
	config        appConfig.Config[*VFSSettings]
	modOnce       sync.Once
}

func (vfs *VFS) ComponentStore() injection.ComponentStore {
	return vfs.store
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
	vfs.modOnce.Do(func() {
		vfs.config = appConfig.NewConfig[*VFSSettings]("extfs.toml")
	})
	return []interface{}{
		vfs.config,
	}
}

func (vfs *VFS) OnConfigUpdated(settings appConfig.AppSettings) {
	vfs.locker.Lock()
	defer vfs.locker.Unlock()
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
