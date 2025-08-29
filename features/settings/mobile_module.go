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
	settingsSrv := &SettingsServlet{}
	mobileSettingsSrv := &MobileSettingsServlet{}
	mobileSettingsSvc := &MobileSettingsService{}

	mobileModule := &stdMobileSettingsModule{}
	mobileModule.module = module
	mobileModule.settingsSrv = settingsSrv
	mobileModule.mobileSettingsSrv = mobileSettingsSrv
	mobileModule.mobileSettingsSvc = mobileSettingsSvc

	module.configListeners = append(module.configListeners, mobileModule)
	module.isMobileMode = true

	return mobileModule
}

const (
	MobileSettingsModuleName = "mobile-settings"
)

type stdMobileSettingsModule struct {
	module *stdModule

	settingsSrv       *SettingsServlet
	mobileSettingsSrv *MobileSettingsServlet
	mobileSettingsSvc *MobileSettingsService
}

var _ = (SettingsConfigListener)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnConfigUpdated(cfg SettingsConfig) {
	if mobileCfg, ok := cfg.(MobileSettingsConfig); ok {
		m.configure(mobileCfg)
	}
}

var _ = (serlvet.SerlvetModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) SetupToSerlvet(app serlvet.SerlvetApp) error {
	err := m.settingsSrv.SetupToSerlvet(app.Route([]byte(SettingsModuleName)))
	if err == nil {
		err = m.mobileSettingsSrv.SetupToSerlvet(app.Route([]byte(MobileSettingsModuleName)))
	}
	return err
}

var _ = (injection.ComponentStoreProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) ComponentStore() injection.ComponentStore {
	return m.module.ComponentStore()
}

var _ = (injection.ComponentProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Components() []injection.Component {
	return []injection.Component{
		// controller
		injection.NewComponent(m.settingsSrv, injection.ComponentNoneScope),
		injection.NewComponent(m.mobileSettingsSrv, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.mobileSettingsSvc, injection.ComponentInternalScope),
	}
}

func (m *stdMobileSettingsModule) configure(cfg MobileSettingsConfig) error {
	mobileSettings, err := m.mobileSettingsSvc.Load()
	if err != nil {
		m.module.logger.Error("MobileSettingsModule", "configure Error: load failed"+err.Error())
		return err
	}

	useWifiOnly := false
	var wifiIfaces []NetInterface
	if mobileSettings.WifiOnly != nil {
		useWifiOnly = *mobileSettings.WifiOnly
	} else {
		wifiIfaces = cfg.WifiInterfaces()
		useWifiOnly = len(wifiIfaces) > 0
	}

	if !useWifiOnly {
		return m.module.configure(nil, false)
	}

	if len(wifiIfaces) > 0 {
		return m.module.configure(wifiIfaces, true)
	}
	return m.module.configure(cfg.WifiInterfaces(), true)

}
