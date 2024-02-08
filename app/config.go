package app

import (
	"os"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
	app   *App
}

func NewConfig(app *App) *Config {
	cfg := new(Config)
	cfg.app = app
	cfg.viper = viper.New()
	return cfg
}

func (cfg *Config) init() error {
	settings := cfg.app.settings
	rootPath := settings.AppRoot()
	_, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rootPath, 0755)
	}
	if err == nil {
		configPath := path.Join(rootPath, settings.AppName()+".toml")
		cfg.viper.SetConfigFile(configPath)
	}

	return err
}

func (cfg *Config) Load() error {
	err := cfg.viper.ReadInConfig()
	if os.IsNotExist(err) {
		return cfg.Save()
	}
	if err == nil {
		err = cfg.viper.Unmarshal(cfg.app.settings)
	}

	return err
}

func (cfg *Config) Save() error {
	return cfg.viper.WriteConfig()
}
