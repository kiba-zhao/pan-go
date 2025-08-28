//go:build android || ios

package settings

import (
	"pan/lib/injection"
	"pan/lib/serlvet"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newMobileModule)
}

func newMobileModule(module *stdModule) interface{} {
	mobileSettingsSrv := &MobileSettingsServlet{}

	mobileModule := &stdMobileSettingsModule{}
	mobileModule.module = module
	mobileModule.mobileSettingsSrv = mobileSettingsSrv

	module.configListeners = append(module.configListeners, mobileModule)

	return mobileModule
}

const (
	MobileSettingsModuleName = "mobile-settings"
)

type stdMobileSettingsModule struct {
	module *stdModule

	mobileSettingsSrv *MobileSettingsServlet
}

var _ = (SettingsConfigListener)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnConfigUpdated(cfg SettingsConfig) {
	// TODO: implement
}

var _ = (serlvet.SerlvetModule)((*stdModule)(nil))

func (m *stdMobileSettingsModule) SetupToSerlvet(app serlvet.SerlvetApp) error {
	err := m.mobileSettingsSrv.SetupToSerlvet(app.Route([]byte(MobileSettingsModuleName)))
	return err
}

var _ = (injection.ComponentStoreProvider)((*stdModule)(nil))

func (m *stdMobileSettingsModule) ComponentStore() injection.ComponentStore {
	return m.module.ComponentStore()
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdMobileSettingsModule) Components() []injection.Component {
	return []injection.Component{
		// controller
		injection.NewComponent(m.mobileSettingsSrv, injection.ComponentInternalScope),
		// service
		injection.NewComponent(&MobileSettingsService{}, injection.ComponentInternalScope),
	}
}
