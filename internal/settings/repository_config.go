package settings

import "pan/internal/repository"

type stdRepositoryConfig struct {
	settingsCfg SettingsConfig
}

func newRepositoryConfig(settingsCfg SettingsConfig) repository.RepositoryConfig {
	cfg := &stdRepositoryConfig{}
	cfg.settingsCfg = settingsCfg
	return cfg
}

var _ = (repository.RepositoryConfig)((*stdRepositoryConfig)(nil))

func (cfg *stdRepositoryConfig) BasePath() string {
	return cfg.settingsCfg.HomePath()
}

func (cfg *stdRepositoryConfig) TempPath() string {
	return cfg.settingsCfg.TempPath()
}
