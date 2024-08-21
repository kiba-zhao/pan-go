package app

import (
	"encoding/base64"
	"pan/app/config"
	"pan/app/controllers"
	"pan/app/net"
	"pan/app/node"
	"pan/app/services"
	"pan/runtime"
	"sync"
)

func New() interface{} {
	return runtime.NewModule(&runtime.Injector{}, config.New(), node.New(), net.New(), NewSample(&module{}))
}

const moduleName = "app"

type module struct {
	Node        node.NodeModule
	Config      config.Config[config.AppSettings]
	DBProvider  RepositoryDBProvider
	settings    config.AppSettings
	settingsRW  sync.RWMutex
	controllers []interface{}
	once        sync.Once
}

func (m *module) Name() string {
	return moduleName
}

func (m *module) Controllers() []interface{} {
	m.once.Do(func() {
		// TODO: add web and node controllers
		m.controllers = []interface{}{
			&controllers.NodeController{},
			&controllers.DiskFileController{},
			&controllers.SettingsController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return nil
}

func (m *module) Components() []runtime.Component {
	// base
	components := []runtime.Component{
		runtime.NewComponent(m.DBProvider, runtime.ComponentInternalScope),
		// services
		runtime.NewComponent(&services.DiskFileService{}, runtime.ComponentInternalScope),
		runtime.NewComponent(&services.SettingsService{Provider: m}, runtime.ComponentInternalScope),
	}

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, runtime.NewComponent(ctrl, runtime.ComponentNoneScope))
	}

	return components
}

func (m *module) OnConfigUpdated(settings config.AppSettings) {
	m.settingsRW.Lock()
	defer m.settingsRW.Unlock()
	m.settings = settings
}

func (m *module) Settings() config.Settings {
	m.settingsRW.RLock()
	defer m.settingsRW.RUnlock()

	settings := *m.settings
	return settings
}

func (m *module) SetSettings(settings config.Settings) error {
	return m.Config.Save(&settings)
}

func (m *module) NodeID() string {
	if m.Node == nil {
		return ""
	}

	nodeSettings := m.Node.NodeSettings()
	if nodeSettings == nil || !nodeSettings.Available() {
		return ""
	}
	return base64.StdEncoding.EncodeToString(nodeSettings.NodeID())
}
