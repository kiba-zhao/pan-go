package settings

import "github.com/spf13/viper"

func initViper(viper *viper.Viper, configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(SettingsModuleName)
	viper.SetConfigType("toml")
}
