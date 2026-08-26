//go:build android || ios

package settings

import (
	"pan/pkg/app"
	"pan/pkg/injection"
	"pan/pkg/module"
	"pan/pkg/ptp"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newMobileModule)
}

func newMobileModule(m *stdModule) interface{} {
	deviceInfoApplet := &DeviceInfoAppletModule{}
	deviceNetworkApplet := &DeviceNetworkAppletModule{}
	mobileNetworkApplet := &MobileNetworkAppletModule{}
	mobileNetworkSvc := &MobileNetworkService{}

	mobileModule := &stdMobileSettingsModule{}
	mobileModule.SubModule = module.NewSubModule(mobileModule, m)
	mobileModule.deviceInfoApplet = deviceInfoApplet
	mobileModule.deviceNetworkApplet = deviceNetworkApplet
	mobileModule.mobileNetworkApplet = mobileNetworkApplet
	mobileModule.mobileNetworkSvc = mobileNetworkSvc

	return mobileModule
}

type stdMobileSettingsModule struct {
	SettingsConfigurer SettingsConfigurer
	SecurityConfigurer SecurityConfigurer

	*module.SubModule[*stdMobileSettingsModule, *stdModule]

	deviceInfoApplet    *DeviceInfoAppletModule
	deviceNetworkApplet *DeviceNetworkAppletModule
	mobileNetworkApplet *MobileNetworkAppletModule
	mobileNetworkSvc    *MobileNetworkService
}

var _ = (app.AppletModule)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) SetupToApplet(router app.AppServletRouter) error {
	router_ := router.Route([]byte(SettingsModuleName))
	err := m.deviceInfoApplet.SetupToApplet(router_)
	if err == nil {
		err = m.deviceNetworkApplet.SetupToApplet(router_)
	}
	if err == nil {
		err = m.mobileNetworkApplet.SetupToApplet(router_)
	}
	return err
}

var _ = (injection.ComponentProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		injection.NewComponent[MobileNetworkChangedTrigger](m, injection.ComponentInternalScope),
		// applet
		injection.NewComponent(m.deviceInfoApplet, injection.ComponentNoneScope),
		injection.NewComponent(m.deviceNetworkApplet, injection.ComponentNoneScope),
		injection.NewComponent(m.mobileNetworkApplet, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.mobileNetworkSvc, injection.ComponentInternalScope),
	}
}

var _ = (MobileNetworkChangedTrigger)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) OnMobileNetworkChanged(mobileNetwork MobileNetwork) {
	m.ParentModule().configure(nil)
}

func (m *stdMobileSettingsModule) loadMobileNetwork() (MobileNetwork, error) {
	mobileNetwork, err := m.mobileNetworkSvc.Load()
	if err != nil {
		m.ParentModule().logger.Error("settings.MobileSettingsModule", "loadMobileNetwork Error: "+err.Error())
	}
	return mobileNetwork, err
}

func (m *stdMobileSettingsModule) loadMobileConfig() MobileSettingsConfig {
	settingsCfg := m.SettingsConfigurer.Config()
	if mobileCfg, ok := settingsCfg.(MobileSettingsConfig); ok {
		return mobileCfg
	}
	return nil
}

func (m *stdMobileSettingsModule) NetInterfaces(mobileNetwork *MobileNetwork) ([]NetInterface, []string, bool) {
	mobileCfg := m.loadMobileConfig()
	if mobileCfg == nil {
		return nil, nil, false
	}

	useWifiOnly := false
	var wifiIfaces []NetInterface
	var zoneList []string
	if mobileNetwork.WifiOnly != nil {
		useWifiOnly = *mobileNetwork.WifiOnly
	} else {
		wifiIfaces = mobileCfg.WifiInterfaces()
		zoneList = mobileCfg.WifiZoneList()
		useWifiOnly = len(wifiIfaces) > 0
	}

	if !useWifiOnly {
		settingsCfg := m.SettingsConfigurer.Config()
		return mobileCfg.NetInterfaces(), settingsCfg.IPv6ZoneList(), true
	}

	if mobileNetwork.WifiOnly == nil {
		return wifiIfaces, zoneList, false
	}

	return mobileCfg.WifiInterfaces(), mobileCfg.WifiZoneList(), false

}

var _ = (ModuleNetProvider)((*stdMobileSettingsModule)(nil))

func (m *stdMobileSettingsModule) NewNetConfig(deviceNetwork *DeviceNetwork) (ptp.QuicConfig, ptp.BroadcastConfig) {
	mobileNetwork, err := m.loadMobileNetwork()
	if err != nil {
		return nil, nil
	}

	ifaces, zoneList, flexible := m.NetInterfaces(&mobileNetwork)
	if len(ifaces) <= 0 && flexible {
		return nil, nil
	}

	addrs := make([]string, 0)
	for _, iface := range ifaces {
		addrs = append(addrs, iface.Addr)
	}
	quicCfg := newQuicConfig(deviceNetwork, m.SecurityConfigurer.Config(), addrs)

	settingsConfig := m.SettingsConfigurer.Config()
	broadcastCfg := &stdBroadcastConfig{}
	broadcastCfg.addrs = deviceNetwork.BroadcastAddrs
	broadcastCfg.ipv6Enabled = settingsConfig.IPv6Enabled()
	broadcastCfg.ipv6ZoneList = zoneList
	broadcastCfg.mtu = settingsConfig.MTU()
	broadcastCfg.ifaces = ifaces

	return quicCfg, broadcastCfg
}
