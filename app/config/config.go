package config

import (
	"io/fs"
	"os"
	"pan/app/injection"
	"pan/runtime"
	"path"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

const PackageName = "pan-go"
const RootPathName = "rootPath"
const DefaultRootName = "." + PackageName

type AppConfig = Config[AppSettings]

type ConfigListener[T any] interface {
	OnConfigUpdated(settings T)
}

type Config[T any] interface {
	SetDefaults(settings T)
	Read() (settings T, err error)
	Load() (settings T, err error)
	Save(settings T) error
	ConfigFilePath() string
}

type configImpl[T any] struct {
	rw         sync.RWMutex
	registry   runtime.Registry
	registryRW sync.RWMutex

	viper     *viper.Viper
	isPtrType bool
}

func NewConfig[T any](name string) (Config[T], error) {

	cfg := &configImpl[T]{}

	t := reflect.TypeFor[T]()
	cfg.isPtrType = t.Kind() == reflect.Ptr
	cfg.viper = viper.New()

	rootPath, err := getConfigRootPath()
	if err != nil {
		return nil, err
	}
	cfg.viper.SetConfigFile(path.Join(rootPath, name))

	err = cfg.EnsureConfig()
	if err != nil {
		return nil, err
	}

	err = cfg.viper.ReadInConfig()
	if _, ok := err.(*fs.PathError); ok || os.IsNotExist(err) {
		err = nil
	}

	return cfg, err
}

func (c *configImpl[T]) Init(registry runtime.Registry) error {

	c.registryRW.Lock()
	c.registry = registry
	c.registryRW.Unlock()

	// init settings
	_, err := c.Load()
	return err

}

func (c *configImpl[T]) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ConfigListener[T]](),
	}
}

func (c *configImpl[T]) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent[Config[T]](c, injection.ComponentExternalScope),
	}
}

func (c *configImpl[T]) SetDefaults(settings T) {
	c.rw.Lock()
	defer c.rw.Unlock()
	t := reflect.TypeOf(settings)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
	v := reflect.ValueOf(settings)
	iv := reflect.Indirect(v)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		fv := iv.FieldByName(field.Name)
		c.viper.SetDefault(field.Name, fv.Interface())
	}
}

func (c *configImpl[T]) Read() (T, error) {
	var settings T
	var err error
	if c.isPtrType {
		settings = reflect.New(reflect.TypeOf(settings).Elem()).Interface().(T)
		err = c.viper.Unmarshal(settings)
	} else {
		err = c.viper.Unmarshal(&settings)
	}

	return settings, err
}

func (c *configImpl[T]) Load() (T, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	settings, err := c.Read()
	if err == nil {
		c.registryRW.RLock()
		registry := c.registry
		c.registryRW.RUnlock()
		onSettingsUpdated(registry, settings)
	}
	return settings, err
}

func (c *configImpl[T]) Save(settings T) error {

	c.rw.Lock()
	needSave := false
	t := reflect.TypeOf(settings)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
	v := reflect.ValueOf(settings)
	iv := reflect.Indirect(v)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}
		fv := iv.FieldByName(field.Name).Interface()
		value := c.viper.Get(field.Name)
		if reflect.DeepEqual(value, fv) {
			continue
		}
		c.viper.Set(field.Name, fv)
		needSave = true
	}
	if !needSave {
		return nil
	}

	err := c.viper.WriteConfig()
	c.rw.Unlock()

	if err == nil {
		_, err = c.Load()
	}
	return err
}

func (c *configImpl[T]) EnsureConfig() error {
	configFilePath := c.ConfigFilePath()
	configDirPath := path.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}

	return err
}

func (c *configImpl[T]) ConfigFilePath() string {
	return c.viper.ConfigFileUsed()
}

func getConfigRootPath() (string, error) {
	rootPath, ok := os.LookupEnv(RootPathName)
	if !ok {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		rootPath = path.Join(homePath, DefaultRootName)
	}

	return rootPath, nil

}

func setDefaultSettings[T any](viper *viper.Viper, settings T) {
	// set defaults
	t := reflect.TypeOf(settings)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fields := reflect.VisibleFields(t)
	v := reflect.ValueOf(settings)
	iv := reflect.Indirect(v)
	for _, field := range fields {
		if !field.IsExported() {
			continue
		}

		fv := iv.FieldByName(field.Name)
		viper.SetDefault(field.Name, fv.Interface())
	}
}

func onSettingsUpdated[T any](registry runtime.Registry, settings T) {
	if registry == nil {
		return
	}
	listeners := runtime.ModulesForType[ConfigListener[T]](registry)
	for _, listener := range listeners {
		listener.OnConfigUpdated(settings)
	}
}

func New() AppConfig {
	settings := newDefaultSettings()
	cfg, err := NewConfig[AppSettings]("pan.toml")
	if err != nil {
		panic(err)
	}
	cfg.SetDefaults(settings)

	return cfg
}
