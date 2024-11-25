package app

import (
	"encoding/base64"
	"path"

	appnode "pan/app/app_node"
	appsettings "pan/app/app_settings"
	"pan/app/bootstrap"
	"pan/app/broadcast"
	"pan/app/config"
	diskfile "pan/app/disk_file"

	"pan/app/guard"
	"pan/app/peer"
	"pan/app/quic"
	"pan/app/sample"
	"pan/app/web"
	"pan/runtime"
	"sync"
)

func New() interface{} {
	m := &module{}
	m.peerGuard = &guard.PeerGuard{}
	return runtime.NewModule(bootstrap.New(), config.New(), peer.New(), broadcast.New(), quic.New(), web.New(), sample.New(m))
}

func Bootstrap() interface{} {
	return bootstrap.Bootstrap()
}

const moduleName = "app"

type module struct {
	PeerModule      peer.PeerModule
	Config          config.AppConfig
	DB              sample.RepositoryDB
	settings        config.AppSettings
	settingsRW      sync.RWMutex
	controllers     []web.WebController
	controllersOnce sync.Once
	peerGuard       *guard.PeerGuard
}

func (m *module) Name() string {
	return moduleName
}

func (m *module) WebControllers() []web.WebController {
	m.controllersOnce.Do(func() {
		// TODO: add web and node controllers
		m.controllers = []web.WebController{
			&appnode.AppNodeController{},
			&diskfile.DiskFileController{},
			&appsettings.AppSettingsController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&appnode.AppNode{},
	}
}

func (m *module) Components() []bootstrap.Component {
	// base
	components := []bootstrap.Component{
		// submodules
		bootstrap.NewComponent(m.peerGuard, bootstrap.ComponentNoneScope),
	}

	// services
	components = sample.AppendSampleComponent(components, &diskfile.DiskFileService{})
	components = sample.AppendSampleExternalComponent[appsettings.AppSettingsExternalService](components, &appsettings.AppSettingsService{Provider: m})
	components = sample.AppendSampleExternalComponent[appnode.AppNodeExternalService](components, &appnode.AppNodeService{})

	// repositories
	components = sample.AppendSampleComponent(components, appnode.NewAppNodeRepository(m.DB))

	// controllers
	for _, ctrl := range m.WebControllers() {
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

func (m *module) PeerID() string {
	if m.PeerModule == nil {
		return ""
	}

	settings := m.PeerModule.PeerSettings()
	if settings == nil || !settings.Available() {
		return ""
	}
	return base64.StdEncoding.EncodeToString(settings.PeerID())
}

func (m *module) Modules() []interface{} {
	return []interface{}{m.peerGuard}
}
