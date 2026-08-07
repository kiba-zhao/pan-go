//go:build !(android || ios)

package gobase

import (
	"net"
	"os"
	"pan/internal/settings"
	"path/filepath"

	"pan/pkg/log"
)

type stdSettings struct {
	logger log.Logger
	root   string
}

func NewSettingsConfig(logger log.Logger, root string) settings.SettingsConfig {
	return &stdSettings{logger: logger, root: root}
}

var _ = (settings.SettingsConfig)((*stdSettings)(nil))

func (s *stdSettings) HomePath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		s.logger.Error("gobase", "Settings.HomePath Error: "+err.Error())
		return ""
	}

	if len(s.root) > 0 {
		return filepath.Join(homePath, s.root)
	}
	return homePath
}

func (s *stdSettings) TempPath() string {
	tempPath := os.TempDir()
	if len(s.root) > 0 {
		return filepath.Join(tempPath, s.root)
	}
	return tempPath
}

func (s *stdSettings) CachePath() string {
	cachePath, err := os.UserCacheDir()
	if err != nil {
		s.logger.Error("gobase", "Settings.CachePath Error: "+err.Error())
		return ""
	}

	if len(s.root) > 0 {
		return filepath.Join(cachePath, s.root)
	}
	return cachePath
}

func (s *stdSettings) HostName() string {
	hostName, err := os.Hostname()
	if err == nil {
		return hostName
	}

	s.logger.Error("gobase", "Settings.HostName Error: "+err.Error())
	return ""
}

func (s *stdSettings) MTU() int {
	mtu := 0
	ifaces, err := net.Interfaces()
	if err != nil || len(ifaces) <= 0 {
		return mtu
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagMulticast == 0 {
			continue
		}
		if iface.MTU < mtu {
			mtu = iface.MTU
		}
	}
	return mtu
}
func (s *stdSettings) IPv6ZoneList() []string {
	return nil
}
func (s *stdSettings) IPv6Enabled() bool {
	enabled := false
	ifaces, err := net.Interfaces()
	if err != nil || len(ifaces) <= 0 {
		return enabled
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagMulticast == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil || len(addrs) <= 0 {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipNet.IP.To4() == nil && !ipNet.IP.IsLoopback() {
				enabled = true
				break
			}
		}
	}
	return enabled
}
