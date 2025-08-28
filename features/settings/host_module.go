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
	hostCtrl := &HostSettingsController{}

	hostModule := &stdHostSettingsModule{}
	hostModule.module = module
	hostModule.hostSettingsCtrl = hostCtrl

	module.configListeners = append(module.configListeners, hostModule)

	return hostModule
}

const (
	HostSettingsModuleName = "host-settings"
)

type stdHostSettingsModule struct {
	module *stdModule

	hostSettingsCtrl *HostSettingsController
}

var _ = (SettingsConfigListener)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) OnConfigUpdated(cfg SettingsConfig) {
	// TODO: implement
}

var _ = (web.WebAppModule)((*stdModule)(nil))

func (m *stdHostSettingsModule) SetupToWeb(app web.WebApp) error {
	err := m.hostSettingsCtrl.SetupToWeb(app.Group(HostSettingsModuleName))
	return err
}

var _ = (injection.ComponentStoreProvider)((*stdModule)(nil))

func (m *stdHostSettingsModule) ComponentStore() injection.ComponentStore {
	return m.module.ComponentStore()
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdHostSettingsModule) Components() []injection.Component {
	return []injection.Component{
		// controller
		injection.NewComponent(m.hostSettingsCtrl, injection.ComponentInternalScope),
		// service
		injection.NewComponent(&HostSettingsService{}, injection.ComponentInternalScope),
	}
}
