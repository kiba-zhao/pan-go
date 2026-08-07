//go:build android || ios

package settings

import (
	"pan/pkg/app"
	"pan/pkg/injection"
	"pan/pkg/module"
	"pan/pkg/net"
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

	// m.ignoreConfigured = true

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

var _ = (app.AppletModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) SetupToApplet(router app.AppServletRouter) error {
	err := m.baseSettingsModule.SetupToApplet(router.Route([]byte(SettingsModuleName)))
	if err == nil {
		err = m.mobileSettingsModule.SetupToApplet(router.Route([]byte(MobileSettingsModuleName)))
	}
	return err
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

var _ = (MobileSettingsChangedTrigger)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnMobileSettingsChanged(mobileSettings MobileSettings) {
	m.ParentModule().configure(nil)
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

func (m *stdMobileSettingsModule) NetInterfaces(mobileSettings *MobileSettings) ([]NetInterface, []string, bool) {
	mobileCfg := m.loadMobileConfig()
	if mobileCfg == nil {
		return nil, nil, false
	}

	useWifiOnly := false
	var wifiIfaces []NetInterface
	var zoneList []string
	if mobileSettings.WifiOnly != nil {
		useWifiOnly = *mobileSettings.WifiOnly
	} else {
		wifiIfaces = mobileCfg.WifiInterfaces()
		zoneList = mobileCfg.WifiZoneList()
		useWifiOnly = len(wifiIfaces) > 0
	}

	if !useWifiOnly {
		settingsCfg := m.SettingsConfigurer.Config()
		return mobileCfg.NetInterfaces(), settingsCfg.IPv6ZoneList(), true
	}

	if mobileSettings.WifiOnly == nil {
		return wifiIfaces, zoneList, false
	}

	return mobileCfg.WifiInterfaces(), mobileCfg.WifiZoneList(), false

}

var _ = (ModuleNetProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) NewNetConfig(settings *Settings) (net.QuicConfig, net.BroadcastConfig) {
	mobileSettings, err := m.loadMobileSettings()
	if err != nil {
		return nil, nil
	}

	ifaces, zoneList, flexible := m.NetInterfaces(&mobileSettings)
	if len(ifaces) <= 0 && flexible {
		return nil, nil
	}

	addrs := make([]string, 0)
	for _, iface := range ifaces {
		addrs = append(addrs, iface.Addr)
	}
	quicCfg := newQuicConfig(settings, m.SecurityConfigurer.Config(), addrs)

	settingsConfig := m.SettingsConfigurer.Config()
	broadcastCfg := &stdBroadcastConfig{}
	broadcastCfg.addrs = settings.BroadcastAddrs
	broadcastCfg.ipv6Enabled = settingsConfig.IPv6Enabled()
	broadcastCfg.ipv6ZoneList = zoneList
	broadcastCfg.mtu = settingsConfig.MTU()
	broadcastCfg.ifaces = ifaces

	return quicCfg, broadcastCfg
}
