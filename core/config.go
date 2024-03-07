package core

import (
	"io/fs"
	"os"
	"path"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

type ConfigModule interface {
	OnInitConfig(cfg Config) error
}

type Config interface {
	Init(settings interface{})
	Load(settings interface{}) error
	Sync(settings interface{}) error
}

type configImpl struct {
	viper  *viper.Viper
	locker sync.Mutex
}

func newConfig(configPath string) Config {
	cfg := &configImpl{}
	cfg.viper = viper.New()
	cfg.viper.SetConfigFile(configPath)
	return cfg
}

func (cfg *configImpl) Init(settings interface{}) {

	cfg.locker.Lock()
	defer cfg.locker.Unlock()
	// Set defaults
	t := reflect.TypeOf(settings).Elem()
	v := reflect.ValueOf(settings)
	iv := reflect.Indirect(v)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := iv.FieldByName(ft.Name)
		cfg.viper.SetDefault(ft.Name, fv.Interface())
	}

	// load config
	err := cfg.viper.ReadInConfig()
	if err == nil {
		err = cfg.viper.Unmarshal(settings)
	}

	if err != nil {
		_, ok := err.(*fs.PathError)
		if !ok && err != os.ErrNotExist {
			panic(err)
		}
	}
}

func (cfg *configImpl) Load(settings interface{}) error {
	return cfg.viper.Unmarshal(settings)
}

func (cfg *configImpl) Sync(settings interface{}) error {

	cfg.locker.Lock()
	defer cfg.locker.Unlock()

	configFilePath := cfg.viper.ConfigFileUsed()
	configDirPath := path.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}
	if err != nil {
		return err
	}

	needSave := false
	t := reflect.TypeOf(settings).Elem()
	v := reflect.ValueOf(settings)
	iv := reflect.Indirect(v)
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := iv.FieldByName(ft.Name).Interface()
		value := cfg.viper.Get(ft.Name)
		if reflect.DeepEqual(value, fv) {
			continue
		}
		cfg.viper.Set(ft.Name, fv)
		needSave = true
	}

	if !needSave {
		return nil
	}
	return cfg.viper.WriteConfig()
}

func generateAppConfigPath(settings *Settings) string {
	return path.Join(settings.AppRoot(), settings.AppName()+"1.toml")
}

func generateConfigPath(settings *Settings, name string) string {
	return path.Join(settings.AppRoot(), "conf.d", name+".toml")
}
