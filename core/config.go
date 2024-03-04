package core

import (
	"io/fs"
	"os"
	"path"
	"reflect"

	"github.com/spf13/viper"
)

type ConfigModule interface {
	Settings() interface{}
	OnInitConfig(cfg Config) error
}

type Config interface {
	Load() error
	Save() error
}

type configImpl struct {
	viper    *viper.Viper
	settings interface{}
}

func (cfg *configImpl) init(configPath string, settings interface{}) error {

	cfg.viper = viper.New()
	cfg.settings = settings

	cfg.viper.SetConfigFile(configPath)

	err := cfg.Load()

	if err != nil {
		_, ok := err.(*fs.PathError)
		if ok || err == os.ErrNotExist {
			err = cfg.Save()
		}
	}

	return err
}

func (cfg *configImpl) Load() error {
	err := cfg.viper.ReadInConfig()
	if err == nil {
		err = cfg.viper.Unmarshal(cfg.settings)
	}

	return err
}

func (cfg *configImpl) Save() error {
	configFilePath := cfg.viper.ConfigFileUsed()
	configDirPath := path.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}
	if err != nil {
		return err
	}

	t := reflect.TypeOf(cfg.settings).Elem()
	v := reflect.ValueOf(cfg.settings)
	iv := reflect.Indirect(v)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := iv.FieldByName(ft.Name)
		cfg.viper.Set(ft.Name, fv.Interface())
	}

	return cfg.viper.WriteConfig()
}
