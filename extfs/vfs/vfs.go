package vfs

import (
	"context"
	"pan/app/bootstrap"
	appConfig "pan/app/config"
	"sync"
)

type VFSFileSystem interface {
	Mount(VFSSettings) error
	Unmount() error
}

type VFS struct {
	VFSFileSystem VFSFileSystem
	locker        sync.Mutex
	vfsSettings   *VFSSettings
	needReload    bool
	reloadChan    chan struct{}
	reloadOnce    sync.Once
	config        appConfig.Config[*VFSSettings]
	modOnce       sync.Once
}

func (vfs *VFS) VFSComponents() []bootstrap.Component {
	components := []bootstrap.Component{
		bootstrap.NewComponent(vfs, bootstrap.ComponentNoneScope),
	}

	// append FUSEFileSystem
	var fuse FUSEFileSystem
	components = append(components,
		bootstrap.NewComponent(&fuse, bootstrap.ComponentNoneScope),
		bootstrap.NewComponent[VFSFileSystem](&fuse, bootstrap.ComponentInternalScope),
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

		vfs.VFSFileSystem.Mount(settings)
	}
	return err
}
