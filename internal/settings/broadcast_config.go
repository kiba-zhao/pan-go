package settings

import (
	"pan/pkg/ptp"
)

type stdBroadcastConfig struct {
	ifaces       []ptp.BroadcastInterface
	mtu          int
	ipv6ZoneList []string
	ipv6Enabled  bool

	security SecurityConfig
	addrs    []string
}

func newBroadcastConfig(settings *Settings, settingsConfig SettingsConfig) ptp.BroadcastConfig {
	cfg := &stdBroadcastConfig{}
	cfg.addrs = settings.BroadcastAddrs
	cfg.ipv6Enabled = settingsConfig.IPv6Enabled()
	cfg.ipv6ZoneList = settingsConfig.IPv6ZoneList()
	cfg.mtu = settingsConfig.MTU()
	cfg.ifaces = []ptp.BroadcastInterface{
		ptp.BroadcastInterface{
			MTU: cfg.mtu,
		},
	}
	return cfg
}

var _ = (ptp.BroadcastConfig)((*stdBroadcastConfig)(nil))

func (cfg *stdBroadcastConfig) Interfaces() []ptp.BroadcastInterface {
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
