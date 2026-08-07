//go:build android || ios

package gomobile

import "pan/internal/settings"

type NetInterface struct {
	settings.NetInterface
}

type SettingsConfig interface {
	HomePath() string
	TempPath() string
	CachePath() string
	HostName() string
	MTU() int
	IPv6Enabled() bool
	IPv6ZoneList() []string

	NetInterfaces() []NetInterface
	WifiInterfaces() []NetInterface
	WifiZoneList() []string
}

type stdSettingsConfig struct {
	cfg SettingsConfig
}

var _ = (settings.SettingsConfig)((*stdSettingsConfig)(nil))

func (s *stdSettingsConfig) HomePath() string {
	return s.cfg.HomePath()
}

func (s *stdSettingsConfig) TempPath() string {
	return s.cfg.TempPath()
}

func (s *stdSettingsConfig) CachePath() string {
	return s.cfg.CachePath()
}

func (s *stdSettingsConfig) HostName() string {
	return s.cfg.HostName()
}

func (s *stdSettingsConfig) MTU() int {
	return s.cfg.MTU()
}

func (s *stdSettingsConfig) IPv6Enabled() bool {
	return s.cfg.IPv6Enabled()
}
func (s *stdSettingsConfig) IPv6ZoneList() []string {
	return s.cfg.IPv6ZoneList()
}

var _ = (settings.MobileSettingsConfig)((*stdSettingsConfig)(nil))

func (s *stdSettingsConfig) NetInterfaces() []settings.NetInterface {
	var netInterfaces []settings.NetInterface
	for _, iface := range s.cfg.NetInterfaces() {
		netInterfaces = append(netInterfaces, iface.NetInterface)
	}
	return netInterfaces
}

var _ = (settings.MobileSettingsConfig)((*stdSettingsConfig)(nil))

func (s *stdSettingsConfig) WifiInterfaces() []settings.NetInterface {
	var netInterfaces []settings.NetInterface
	for _, iface := range s.cfg.WifiInterfaces() {
		netInterfaces = append(netInterfaces, iface.NetInterface)
	}
	return netInterfaces
}

func (s *stdSettingsConfig) WifiZoneList() []string {
	return s.cfg.WifiZoneList()
}
