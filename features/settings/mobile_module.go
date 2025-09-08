//go:build android || ios

package settings

import (
	"context"
	"pan/lib/bootstrap"
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

	module.ignoreConfigured = true

	return mobileModule
}

const (
	MobileSettingsModuleName = "mobile-settings"
)

type stdMobileSettingsModule struct {
	SettingsConfigurer SettingsConfigurer
	SecurityConfigurer SecurityConfigurer

	module *stdModule

	settingsSrv       *SettingsServlet
	mobileSettingsSrv *MobileSettingsServlet
	mobileSettingsSvc *MobileSettingsService
}

var _ = (SecurityConfigurerListener)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnConfigUpdated(cfg SecurityConfig) {
	mobileCfg, ok := cfg.(MobileSettingsConfig)
	if !ok {
		return
	}

	var settings Settings
	mobileSettings, err := m.loadMobileSettings()
	if err == nil {
		settings, err = m.module.loadSettings()
	}

	if err != nil {
		return
	}

	m.configure(cfg, settings, mobileSettings, mobileCfg)
}

var _ = (serlvet.SerlvetModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) SetupToSerlvet(app serlvet.SerlvetApp) error {
	err := m.settingsSrv.SetupToSerlvet(app.Route([]byte(SettingsModuleName)))
	if err == nil {
		err = m.mobileSettingsSrv.SetupToSerlvet(app.Route([]byte(MobileSettingsModuleName)))
	}
	return err
}

var _ = (bootstrap.DeferModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Defer(ctx context.Context) error {
	m.SecurityConfigurer.Subscribe(m)
	return nil
}

var _ = (bootstrap.DestroyModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Destroy() {
	m.SecurityConfigurer.Unsubscribe(m)
}

var _ = (injection.ComponentStoreProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) ComponentStore() injection.ComponentStore {
	return m.module.ComponentStore()
}

var _ = (injection.ComponentProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		injection.NewComponent[MobileSettingsChangedTrigger](m, injection.ComponentInternalScope),
		// controller
		injection.NewComponent(m.settingsSrv, injection.ComponentNoneScope),
		injection.NewComponent(m.mobileSettingsSrv, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.mobileSettingsSvc, injection.ComponentInternalScope),
	}
}

var _ = (SettingsChangedTrigger)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnSettingsChanged(settings Settings) {

	mobileSettings, err := m.loadMobileSettings()
	if err != nil {
		return
	}

	mobileCfg := m.loadMobileConfig()
	if mobileCfg == nil {
		return
	}

	securityCfg := m.SecurityConfigurer.Config()
	m.configure(securityCfg, settings, mobileSettings, mobileCfg)
}

var _ = (MobileSettingsChangedTrigger)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnMobileSettingsChanged(mobileSettings MobileSettings) {
	settings, err := m.module.loadSettings()
	if err != nil {
		return
	}
	mobileCfg := m.loadMobileConfig()
	if mobileCfg == nil {
		return
	}

	securityCfg := m.SecurityConfigurer.Config()
	m.configure(securityCfg, settings, mobileSettings, mobileCfg)
}

func (m *stdMobileSettingsModule) configure(securityCfg SecurityConfig, settings Settings, mobileSettings MobileSettings, cfg MobileSettingsConfig) error {

	useWifiOnly := false
	var wifiIfaces []NetInterface
	if mobileSettings.WifiOnly != nil {
		useWifiOnly = *mobileSettings.WifiOnly
	} else {
		wifiIfaces = cfg.WifiInterfaces()
		useWifiOnly = len(wifiIfaces) > 0
	}

	if !useWifiOnly {
		return m.module.configure(securityCfg, settings, nil, false)
	}

	if len(wifiIfaces) > 0 {
		return m.module.configure(securityCfg, settings, wifiIfaces, true)
	}
	return m.module.configure(securityCfg, settings, cfg.WifiInterfaces(), true)

}

func (m *stdMobileSettingsModule) loadMobileSettings() (MobileSettings, error) {
	mobileSettings, err := m.mobileSettingsSvc.Load()
	if err != nil {
		m.module.logger.Error("MobileSettingsModule", "loadMobileSettings Error: "+err.Error())
	}
	return mobileSettings, err
}

func (m *stdMobileSettingsModule) loadMobileConfig() MobileSettingsConfig {
	settingsCfg := m.SettingsConfigurer.Config()
	if mobileCfg, ok := settingsCfg.(MobileSettingsConfig); ok {
		return mobileCfg
	}
	return nil
}
