package app

import (
	"encoding/base64"
	"path"

	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/constant"
	"pan/app/controllers"
	"pan/app/models"
	"pan/app/net"
	"pan/app/node"
	"pan/app/repositories"
	repoImpl "pan/app/repositories/impl"
	"pan/app/services"
	"pan/runtime"
	"sync"
)

func New() interface{} {
	m := &module{}
	m.guard = &guardModule{}
	return runtime.NewModule(bootstrap.New(), config.New(), node.New(), net.New(), NewSample(m))
}

func Bootstrap() interface{} {
	return bootstrap.Bootstrap()
}

const moduleName = "app"

type module struct {
	Node        node.NodeModule
	Config      config.AppConfig
	DBProvider  RepositoryDBProvider
	settings    config.AppSettings
	settingsRW  sync.RWMutex
	controllers []interface{}
	once        sync.Once
	guard       *guardModule
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
	return []interface{}{
		&models.Node{},
	}
}

func (m *module) Components() []bootstrap.Component {
	// base
	components := []bootstrap.Component{
		bootstrap.NewComponent(m.DBProvider, bootstrap.ComponentInternalScope),
		// submodules
		bootstrap.NewComponent(m.guard, bootstrap.ComponentNoneScope),
	}

	// services
	components = AppendSampleComponent(components, &services.DiskFileService{})
	components = AppendSampleExternalComponent[services.SettingsExternalService](components, &services.SettingsService{Provider: m})
	components = AppendSampleExternalComponent[services.NodeExternalService](components, &services.NodeService{Provider: m})

	// repositories
	components = AppendSampleComponent[repositories.NodeRepository](components, &repoImpl.NodeRepository{})

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, bootstrap.NewComponent(ctrl, bootstrap.ComponentNoneScope))
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

func (m *module) RootPath() string {
	return path.Dir(m.Config.ConfigFilePath())
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

func (m *module) NodeManager() node.NodeManager {
	if m.Node == nil {
		return nil
	}
	mgr := m.Node.NodeManager()
	return mgr
}

func (m *module) Modules() []interface{} {
	return []interface{}{m.guard}
}

type guardModule struct {
	NodeService     *services.NodeService
	SettingsService *services.SettingsService
}

func (g *guardModule) Enabled() bool {
	settings := g.SettingsService.Load()
	return settings.GuardEnabled
}

func (g *guardModule) Access(nodeId node.NodeID) error {
	err := g.NodeService.AccessWithNodeID(nodeId)
	if err != nil {
		return err
	}

	settings := g.SettingsService.Load()
	if !settings.GuardAccess {
		err = constant.ErrRefused
	}
	return err
}
