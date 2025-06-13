package vfs

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/runtime"
)

type VFSConfig = config.Config[*VFSSettings]

type stdVFSModule struct {
	injection.ComponentStoreProvider

	AppConfig config.AppConfig
	VFSConfig VFSConfig

	vfsfsRuntime   *stdVFSFSRuntime
	vfsServer      *stdVFSServer
	vfsConfigProxy *stdVFSConfigProxy
}

func New(provider injection.ComponentStoreProvider) interface{} {

	module := &stdVFSModule{}
	module.ComponentStoreProvider = provider

	vfsfsRuntime := &stdVFSFSRuntime{}
	module.vfsfsRuntime = vfsfsRuntime

	vfsServer := &stdVFSServer{}
	module.vfsServer = vfsServer
	vfsServer.runtime = vfsfsRuntime
	vfsServer.logger = log.Default()
	vfsServer.reloadChan = make(chan struct{}, 1)

	vfsConfigProxy := &stdVFSConfigProxy{}
	module.vfsConfigProxy = vfsConfigProxy
	vfsConfigProxy.module = module

	return module
}

var _ = (injection.ComponentProvider)((*stdVFSModule)(nil))

func (m *stdVFSModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
		injection.NewComponent(m.vfsfsRuntime, injection.ComponentNoneScope),
	}
}

var _ = (config.AppConfigListener)((*stdVFSModule)(nil))

func (m *stdVFSModule) OnConfigUpdated(settings config.AppSettings) {

	defaultsSettings := newDefaultsVFSSettings(settings)
	m.VFSConfig.SetDefaults(defaultsSettings)
	if _, already := m.VFSConfig.SettingsAndAlready(); already {
		m.VFSConfig.Load()
	}

}

var _ = (runtime.ProviderModule)((*stdVFSModule)(nil))

func (m *stdVFSModule) Modules() []interface{} {
	return []interface{}{
		config.New[*VFSSettings]("extfs_vfs.toml"),
	}
}

var _ = (bootstrap.ReadyModule)((*stdVFSModule)(nil))

func (m *stdVFSModule) Ready(ctx context.Context) error {

	m.AppConfig.Subscribe(m)
	defer m.AppConfig.Unsubscribe(m)

	m.VFSConfig.Subscribe(m.vfsConfigProxy)
	defer m.VFSConfig.Unsubscribe(m.vfsConfigProxy)

	return m.vfsServer.RunAndServe(ctx)
}

type stdVFSConfigProxy struct {
	module *stdVFSModule
}

var _ = (config.ConfigListener[*VFSSettings])((*stdVFSConfigProxy)(nil))

func (m *stdVFSConfigProxy) OnConfigUpdated(settings *VFSSettings) {
	module := m.module
	server := module.vfsServer
	server.SetHostName(settings.HostName)
	server.SetMountPath(settings.MountPath)
}
