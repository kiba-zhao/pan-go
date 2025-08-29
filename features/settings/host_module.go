//go:build !(android || ios) || host

package settings

import (
	"pan/lib/injection"
	"pan/lib/web"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(module *stdModule) interface{} {
	settingsCtrl := &SettingsController{}
	hostCtrl := &HostSettingsController{}
	hostSettingsSvc := &HostSettingsService{}

	hostModule := &stdHostSettingsModule{}
	hostModule.module = module
	hostModule.settingsCtrl = settingsCtrl
	hostModule.hostSettingsCtrl = hostCtrl
	hostModule.hostSettingsSvc = hostSettingsSvc

	module.configListeners = append(module.configListeners, hostModule)

	return hostModule
}

const (
	HostSettingsModuleName = "host-settings"
)

type stdHostSettingsModule struct {
	WebConfigurer web.WebConfigurer

	module *stdModule

	settingsCtrl     *SettingsController
	hostSettingsCtrl *HostSettingsController
	hostSettingsSvc  *HostSettingsService
}

var _ = (SettingsConfigListener)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) OnConfigUpdated(cfg SettingsConfig) {
	m.configure()
}

var _ = (web.WebAppModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) SetupToWeb(app web.WebApp) error {
	err := m.settingsCtrl.SetupToWeb(app.Group(SettingsModuleName))
	if err == nil {
		err = m.hostSettingsCtrl.SetupToWeb(app.Group(HostSettingsModuleName))
	}
	return err
}

var _ = (injection.ComponentStoreProvider)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) ComponentStore() injection.ComponentStore {
	return m.module.ComponentStore()
}

var _ = (injection.ComponentProvider)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		// controller
		injection.NewComponent(m.settingsCtrl, injection.ComponentNoneScope),
		injection.NewComponent(m.hostSettingsCtrl, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.hostSettingsSvc, injection.ComponentInternalScope),
	}
}

func (m *stdHostSettingsModule) configure() error {
	hostSettings, err := m.hostSettingsSvc.Load()
	if err != nil {
		m.module.logger.Error("HostSettingsModule", "configure Error: load failed"+err.Error())
		return err
	}
	err = m.configureWeb(&hostSettings)
	return err
}

func (m *stdHostSettingsModule) configureWeb(hostSettings *HostSettings) error {
	webConfig := newWebConfig(hostSettings)
	return m.WebConfigurer.Configure(webConfig)
}
