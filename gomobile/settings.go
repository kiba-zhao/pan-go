//go:build android || ios

package gomobile

import "pan/features/settings"

type NetInterface struct {
	settings.NetInterface
}

type SettingsConfig interface {
	HomePath() string
	TempPath() string
	CachePath() string
	HostName() string
	NetInterfaces() []NetInterface
	WifiInterfaces() []NetInterface
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

func (s *stdSettingsConfig) NetInterfaces() []settings.NetInterface {
	var netInterfaces []settings.NetInterface
	for _, netInterface := range s.cfg.NetInterfaces() {
		netInterfaces = append(netInterfaces, netInterface.NetInterface)
	}
	return netInterfaces
}

var _ = (settings.MobileSettingsConfig)((*stdSettingsConfig)(nil))

func (s *stdSettingsConfig) WifiInterfaces() []settings.NetInterface {
	var netInterfaces []settings.NetInterface
	for _, netInterface := range s.cfg.WifiInterfaces() {
		netInterfaces = append(netInterfaces, netInterface.NetInterface)
	}
	return netInterfaces
}
