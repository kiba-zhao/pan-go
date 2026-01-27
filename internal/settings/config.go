package settings

import "pan/internal/config"

type SettingsConfig interface {
	HomePath() string
	TempPath() string
	CachePath() string

	HostName() string
	NetInterfaces() []NetInterface
}

type SettingsConfigListener = config.ConfigurerListener[SettingsConfig]
type SettingsConfigurer = config.Configurer[SettingsConfig]
