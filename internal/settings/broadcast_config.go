package settings

import (
	"pan/pkg/net"
)

type stdBroadcastConfig struct {
	ifaces       []net.BroadcastInterface
	mtu          int
	ipv6ZoneList []string
	ipv6Enabled  bool

	security SecurityConfig
	addrs    []string
}

func newBroadcastConfig(settings *Settings, settingsConfig SettingsConfig) net.BroadcastConfig {
	cfg := &stdBroadcastConfig{}
	cfg.addrs = settings.BroadcastAddrs
	cfg.ipv6Enabled = settingsConfig.IPv6Enabled()
	cfg.ipv6ZoneList = settingsConfig.IPv6ZoneList()
	cfg.mtu = settingsConfig.MTU()
	cfg.ifaces = []net.BroadcastInterface{
		net.BroadcastInterface{
			MTU: cfg.mtu,
		},
	}
	return cfg
}

// func newBroadcastConfig(settings *Settings, settingsConfig SettingsConfig, netIfaces []NetInterface) net.BroadcastConfig {
// 	cfg := &stdBroadcastConfig{}
// 	cfg.addrs = settings.BroadcastAddrs

// 	if len(netIfaces) <= 0 {
// 		cfg.ipv6Enabled = settingsConfig.IPv6Enabled()
// 		cfg.ipv6ZoneList = settingsConfig.IPv6ZoneList()
// 		cfg.mtu = settingsConfig.MTU()
// 		return cfg
// 	}

// 	cfg.ifaces = make([]net.BroadcastInterface, 0)
// 	cfg.ipv6ZoneList = make([]string, 0)
// 	cfg.mtu = 0
// 	for _, iface := range netIfaces {

// 		var bIface net.BroadcastInterface
// 		bIface.Addr = iface.Addr
// 		bIface.MTU = iface.MTU
// 		cfg.ifaces = append(cfg.ifaces, bIface)

// 		if len(iface.Zone) > 0 {
// 			cfg.ipv6ZoneList = append(cfg.ipv6ZoneList, iface.Zone)
// 		}

// 		if iface.MTU < cfg.mtu {
// 			cfg.mtu = iface.MTU
// 		}
// 	}

// 	return cfg
// }

var _ = (net.BroadcastConfig)((*stdBroadcastConfig)(nil))

func (cfg *stdBroadcastConfig) Interfaces() []net.BroadcastInterface {
	return cfg.ifaces
}

func (cfg *stdBroadcastConfig) Addrs() []string {
	return cfg.addrs
}

func (cfg *stdBroadcastConfig) MTU() int {
	return cfg.mtu
}

func (cfg *stdBroadcastConfig) IPv6ZoneList() []string {
	return cfg.ipv6ZoneList
}

func (cfg *stdBroadcastConfig) IPv6Enabled() bool {
	return cfg.ipv6Enabled
}

func (cfg *stdBroadcastConfig) DeliverMaxSize() int {
	return 65535
}
