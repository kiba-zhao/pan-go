package settings

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func initViper(viper *viper.Viper, configPath string) error {
	viper.SetConfigFile(filepath.Join(configPath, SettingsModuleName+".toml"))
	return viper.ReadInConfig()
}
