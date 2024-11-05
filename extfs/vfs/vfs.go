package vfs

import (
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
	settingsRW    sync.RWMutex
	vfsSettings   *VFSSettings
	hasSig        bool
	sigLocker     sync.Mutex
	sigChan       chan bool
	sigOnce       sync.Once
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
	vfs.settingsRW.Lock()
	defer vfs.settingsRW.Unlock()
	defaultsSettings := newDefaultsVFSSettings(settings)
	vfs.config.SetDefaults(defaultsSettings)
	vfsSettings, err := vfs.config.Load()
	if err != nil {
		panic(err)
	}
	vfs.vfsSettings = vfsSettings
	if vfs.hasSig {
		return
	}
	vfs.hasSig = true
	vfs.SigChan() <- true
}

func (vfs *VFS) SigChan() chan bool {
	vfs.sigOnce.Do(func() {
		vfs.sigChan = make(chan bool, 1)
	})
	return vfs.sigChan
}

func (vfs *VFS) Ready() error {

	for {
		sig := <-vfs.SigChan()
		vfs.sigLocker.Lock()
		vfs.hasSig = false
		settings := *vfs.vfsSettings
		vfs.sigLocker.Unlock()

		vfs.VFSFileSystem.Unmount()
		if !sig {
			break
		}

		vfs.VFSFileSystem.Mount(settings)
	}
	return nil
}
