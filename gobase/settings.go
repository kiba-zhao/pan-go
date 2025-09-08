//go:build !(android || ios)

package gobase

import (
	"os"
	"pan/features/settings"

	"pan/lib/log"
)

type stdSettings struct {
	logger log.Logger
}

func NewSettingsConfig(logger log.Logger) settings.SettingsConfig {
	return &stdSettings{logger: logger}
}

var _ = (settings.SettingsConfig)((*stdSettings)(nil))

func (s *stdSettings) HomePath() string {
	homePath, err := os.UserHomeDir()
	if err == nil {
		return homePath
	}

	s.logger.Error("gobase", "Settings.HomePath Error: "+err.Error())
	return ""
}

func (s *stdSettings) TempPath() string {
	tempPath := os.TempDir()
	return tempPath
}

func (s *stdSettings) CachePath() string {
	cachePath, err := os.UserCacheDir()
	if err == nil {
		return cachePath
	}

	s.logger.Error("gobase", "Settings.CachePath Error: "+err.Error())
	return ""
}

func (s *stdSettings) HostName() string {
	hostName, err := os.Hostname()
	if err == nil {
		return hostName
	}

	s.logger.Error("gobase", "Settings.HostName Error: "+err.Error())
	return ""
}

func (s *stdSettings) NetInterfaces() []settings.NetInterface {
	return nil
}
