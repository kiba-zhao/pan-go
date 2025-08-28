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

	hostModule := &stdHostSettingsModule{}
	hostModule.module = module
	hostModule.settingsCtrl = settingsCtrl
	hostModule.hostSettingsCtrl = hostCtrl

	module.configListeners = append(module.configListeners, hostModule)

	return hostModule
}

const (
	HostSettingsModuleName = "host-settings"
)

type stdHostSettingsModule struct {
	module *stdModule

	settingsCtrl     *SettingsController
	hostSettingsCtrl *HostSettingsController
}

var _ = (SettingsConfigListener)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) OnConfigUpdated(cfg SettingsConfig) {
	// TODO: implement
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
		// controller
		injection.NewComponent(m.settingsCtrl, injection.ComponentNoneScope),
		injection.NewComponent(m.hostSettingsCtrl, injection.ComponentNoneScope),
		// service
		injection.NewComponent(&HostSettingsService{}, injection.ComponentInternalScope),
	}
}
