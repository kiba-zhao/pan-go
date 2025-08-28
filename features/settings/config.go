package settings

import "pan/lib/config"

type SettingsConfig interface {
	HomePath() string
	TempPath() string
	CachePath() string

	HostName() string
}

type SettingsConfigListener = config.ConfigurerListener[SettingsConfig]
type SettingsConfigurer = config.Configurer[SettingsConfig]
