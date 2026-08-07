package settings

import "pan/pkg/config"

type SettingsConfig interface {
	HomePath() string
	TempPath() string
	CachePath() string

	HostName() string

	MTU() int
	IPv6ZoneList() []string
	IPv6Enabled() bool
}

type SettingsConfigListener = config.ConfigurerListener[SettingsConfig]
type SettingsConfigurer = config.Configurer[SettingsConfig]
