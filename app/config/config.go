// Package config provides the configuration engine
//
// The configuration engine is used to manage the configuration of the application
package config

import (
	"io/fs"
	"os"
	"pan/app/injection"
	"pan/runtime"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

const PackageName = "pan-go"
const RootPathName = "rootPath"
const DefaultRootName = "." + PackageName

type AppConfig = Config[AppSettings]

// ConfigListener is a listener for config updates
// It is called when the config is updated
type ConfigListener[T any] interface {
	// OnConfigUpdated is called when the config is updated.
	//
	// The function is called with the new config settings.
	OnConfigUpdated(settings T)
}

// Config is the configuration engine
type Config[T any] interface {
	// SetDefaults sets the default values of the configuration
	//
	// The function sets the default values of the configuration based on the given settings.
	// The settings are expected to be a struct with fields that are tagged with the
	// "default" tag. The value of the tag is the default value of the field.
	//
	// The function is called with the settings as an argument.
	// It does not return an error.
	SetDefaults(settings T)
	// Read  the configuration settings from the underlying storage.
	//
	// It returns the settings as a value of type T and an error if any occurs
	// during the reading process. If the configuration is stored as a pointer
	// type, the settings will be unmarshalled into a new instance. If it is stored
	// as a value type, the settings will be unmarshalled directly.
	// An error is returned if the unmarshalling process fails.
	Read() (settings T, err error)
	// Load loads the configuration settings from the underlying storage and
	// notifies all registered config listeners about the updates.
	//
	// It returns the settings as a value of type T and an error if any occurs
	// during the loading process. If the configuration is stored as a pointer
	// type, the settings will be unmarshalled into a new instance. If it is stored
	// as a value type, the settings will be unmarshalled directly.
	// An error is returned if the unmarshalling process fails.
	Load() (settings T, err error)
	// Save the given configuration settings to the underlying storage.
	//
	// It marshals the given settings to JSON and saves them to the file specified
	// by ConfigFilePath. An error is returned if the marshalling or saving process fails.
	Save(settings T) error
	// ConfigFilePath returns the path to the configuration file.
	//
	// The path is determined by the configuration name passed to NewConfig and
	// the root path of the application, which can be set using the environment
	ConfigFilePath() string
}

type configImpl[T any] struct {
	rw         sync.RWMutex
	registry   runtime.Registry
	registryRW sync.RWMutex

	viper     *viper.Viper
	isPtrType bool
}

// NewConfig creates a new configuration engine with the given name.
//
// The name parameter determines the filename of the configuration file, which
// is stored in the root path of the application. The root path is determined by
// the environment variable "PAN_ROOT_PATH" or the default root path if the
// environment variable is not set.
//
// The function returns the configuration engine and an error if any occurs
// during the creation process. If the configuration file does not exist, the
// function will create it with default values and return nil for the error.
//
// The configuration engine can be used to read and write configuration settings.
// The settings are stored as a value of type T. If T is a pointer type, the
// settings will be marshalled into a new instance. If T is a value type, the
// settings will be marshalled directly.
func NewConfig[T any](name string) (Config[T], error) {

	cfg := &configImpl[T]{}

	t := reflect.TypeFor[T]()
	cfg.isPtrType = t.Kind() == reflect.Ptr
	cfg.viper = viper.New()

	rootPath, err := getConfigRootPath()
	if err != nil {
		return nil, err
	}
	cfg.viper.SetConfigFile(filepath.Join(rootPath, name))

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

// EnsureConfig ensures that the configuration directory exists.
// It creates the directory if it does not exist.
// It returns an error if any error occurs during the creation process.
func (c *configImpl[T]) EnsureConfig() error {
	configFilePath := c.ConfigFilePath()
	configDirPath := filepath.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}

	return err
}

// ConfigFilePath returns the path to the configuration file.
//
// The path is determined by the configuration name passed to NewConfig and
// the root path of the application, which can be set using the environment
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
		rootPath = filepath.Join(homePath, DefaultRootName)
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

// New returns a new instance of AppConfig. It will load the configuration from the
// default configuration file path and set the default values based on the given
// settings. If the configuration file does not exist, it will be created with the
// default values. If any error occurs during the creation of the AppConfig, it
// will panic.
func New() AppConfig {
	settings := newDefaultSettings()
	cfg, err := NewConfig[AppSettings]("pan.toml")
	if err != nil {
		panic(err)
	}
	cfg.SetDefaults(settings)

	return cfg
}
