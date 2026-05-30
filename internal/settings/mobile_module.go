//go:build android || ios

package settings

import (
	"context"
	"pan/internal/app"
	"pan/internal/bootstrap"
	"pan/internal/injection"
	"pan/internal/module"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newMobileModule)
}

func newMobileModule(m *stdModule) interface{} {
	baseSettingsModule := &BaseSettingsAppletModule{}
	mobileSettingsModule := &MobileSettingsAppletModule{}
	mobileSettingsSvc := &MobileSettingsService{}

	mobileModule := &stdMobileSettingsModule{}
	mobileModule.SubModule = module.NewSubModule(mobileModule, m)
	mobileModule.baseSettingsModule = baseSettingsModule
	mobileModule.mobileSettingsModule = mobileSettingsModule
	mobileModule.mobileSettingsSvc = mobileSettingsSvc

	m.ignoreConfigured = true

	return mobileModule
}

const (
	MobileSettingsModuleName = "mobile-settings"
)

type stdMobileSettingsModule struct {
	SettingsConfigurer SettingsConfigurer
	SecurityConfigurer SecurityConfigurer

	*module.SubModule[*stdMobileSettingsModule, *stdModule]

	baseSettingsModule   *BaseSettingsAppletModule
	mobileSettingsModule *MobileSettingsAppletModule
	mobileSettingsSvc    *MobileSettingsService
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
		settings, err = m.ParentModule().loadSettings()
	}

	if err != nil {
		return
	}

	m.configure(cfg, settings, mobileSettings, mobileCfg)
}

var _ = (app.AppletModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) SetupToApplet(router app.AppServletRouter) error {
	err := m.baseSettingsModule.SetupToApplet(router.Route([]byte(SettingsModuleName)))
	if err == nil {
		err = m.mobileSettingsModule.SetupToApplet(router.Route([]byte(MobileSettingsModuleName)))
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

var _ = (injection.ComponentProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		injection.NewComponent[MobileSettingsChangedTrigger](m, injection.ComponentInternalScope),
		// controller
		injection.NewComponent(m.baseSettingsModule, injection.ComponentNoneScope),
		injection.NewComponent(m.mobileSettingsModule, injection.ComponentNoneScope),
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
	settings, err := m.ParentModule().loadSettings()
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
		return m.ParentModule().configure(securityCfg, settings, nil, false)
	}

	if len(wifiIfaces) > 0 {
		return m.ParentModule().configure(securityCfg, settings, wifiIfaces, true)
	}
	return m.ParentModule().configure(securityCfg, settings, cfg.WifiInterfaces(), true)

}

func (m *stdMobileSettingsModule) loadMobileSettings() (MobileSettings, error) {
	mobileSettings, err := m.mobileSettingsSvc.Load()
	if err != nil {
		m.ParentModule().logger.Error("settings.MobileSettingsModule", "loadMobileSettings Error: "+err.Error())
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
