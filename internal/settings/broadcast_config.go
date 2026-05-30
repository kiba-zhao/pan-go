package settings

import (
	"pan/internal/net"
)

type stdBroadcastConfig struct {
	ifaces       []net.BroadcastInterface
	mtu          int
	ipv6ZoneList []string

	security SecurityConfig
	addrs    []string
}

func newBroadcastConfig(settings *Settings, netIfaces []NetInterface, isSubNet bool) net.BroadcastConfig {
	cfg := &stdBroadcastConfig{}
	cfg.addrs = settings.BroadcastAddrs

	initBroadcastConfigWithNetInterfaces(cfg, netIfaces, isSubNet)
	return cfg
}

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

func (cfg *stdBroadcastConfig) DeliverMaxSize() int {
	return 65535
}

func initBroadcastConfigWithNetInterfaces(cfg *stdBroadcastConfig, netIfaces []NetInterface, isSubNet bool) {

	if len(netIfaces) <= 0 {
		return
	}

	minMTU := netIfaces[0].MTU
	maxMTU := minMTU
	for _, iface := range netIfaces {

		if isSubNet {
			var bIface net.BroadcastInterface
			bIface.Addr = iface.Addr
			bIface.MTU = iface.MTU
			cfg.ifaces = append(cfg.ifaces, bIface)
		}

		if len(iface.Zone) > 0 {
			cfg.ipv6ZoneList = append(cfg.ipv6ZoneList, iface.Zone)
		}

		if iface.MTU < minMTU {
			minMTU = iface.MTU
			continue
		}
		if iface.MTU > maxMTU {
			maxMTU = iface.MTU
			continue
		}
	}

	cfg.mtu = maxMTU
	if len(cfg.ifaces) < 0 {
		cfg.ifaces = []net.BroadcastInterface{
			{
				MTU: minMTU,
			},
		}
	}
}
