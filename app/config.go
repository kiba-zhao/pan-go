package app

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
	OnConfigChanged(settings T, cfg Config[T]) error
}

type Config[T any] interface {
	Read() (settings T, err error)
	Write(settings T) error
}

type configImpl[T any] struct {
	rw        sync.RWMutex
	listeners []ConfigListener[T]
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

	// init settings
	settings, err := c.Read()
	if err != nil {
		return err
	}

	// init listeners
	c.listeners = make([]ConfigListener[T], 0)
	return runtime.TraverseModules(registry, func(listener ConfigListener[T]) error {
		listenerErr := listener.OnConfigChanged(settings, c)
		if listenerErr == nil {
			c.listeners = append(c.listeners, listener)
		}
		return listenerErr
	})
}

func (c *configImpl[T]) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[ConfigListener[T]](),
	}
}

func (c *configImpl[T]) GetComponents() []Component {
	return []Component{
		NewComponent[Config[T]](c, ComponentExternalScope),
	}
}

func (c *configImpl[T]) readSettings() (T, error) {
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

func (c *configImpl[T]) Read() (T, error) {
	c.rw.RLock()
	defer c.rw.RUnlock()

	return c.readSettings()
}

func (c *configImpl[T]) Write(settings T) error {

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
		for _, listener := range c.listeners {
			err = listener.OnConfigChanged(settings, c)
		}
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
