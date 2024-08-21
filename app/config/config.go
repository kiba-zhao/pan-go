package config

import (
	"os"
	"pan/runtime"
	"path"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

type AppConfig = Config[AppSettings]

type ParseConfigPath[T any] func(settings T) string

type ConfigListener[T any] interface {
	OnConfigUpdated(settings T)
}

type Config[T any] interface {
	Read() (settings T, err error)
	Load() (settings T, err error)
	Save(settings T) error
}

type configImpl[T any] struct {
	rw        sync.RWMutex
	registry  runtime.Registry
	viper     *viper.Viper
	isPtrType bool
}

func NewConfig[T any](settings T, parse ParseConfigPath[T]) Config[T] {

	// TODO: check T is a pointer

	t := reflect.TypeOf(settings)
	isPtrType := t.Kind() == reflect.Ptr

	viper := viper.New()
	configPath := parse(settings)
	viper.SetConfigFile(configPath)
	setDefaultSettings(viper, settings)

	return &configImpl[T]{viper: viper, isPtrType: isPtrType}
}

func (c *configImpl[T]) Init(registry runtime.Registry) error {
	c.registry = registry

	// init settings
	_, err := c.Load()
	return err

}

func (c *configImpl[T]) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ConfigListener[T]](),
	}
}

func (c *configImpl[T]) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent[Config[T]](c, runtime.ComponentExternalScope),
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
		onSettingsUpdated(c.registry, settings)
	}
	return settings, err
}

func (c *configImpl[T]) Save(settings T) error {

	c.rw.Lock()
	defer c.rw.Unlock()

	configFilePath := c.viper.ConfigFileUsed()
	configDirPath := path.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}
	if err != nil {
		return err
	}

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

	err = c.viper.WriteConfig()
	if err == nil {
		onSettingsUpdated(c.registry, settings)
	}
	return err
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
	listeners := runtime.ModulesForType[ConfigListener[T]](registry)
	for _, listener := range listeners {
		listener.OnConfigUpdated(settings)
	}
}

func parseDefaultConfigPath(settings AppSettings) string {
	return path.Join(settings.RootPath, "pan.toml")
}

func New() AppConfig {
	settings := newDefaultSettings()
	config := NewConfig(settings, parseDefaultConfigPath)
	return config
}
