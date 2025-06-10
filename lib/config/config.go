// Package config provides the configuration engine
//
// The configuration engine is used to manage the configuration of the application
package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

// ConfigListener is a listener for config updates
// It is called when the config is updated
type ConfigListener[T any] interface {
	// OnConfigUpdated is called when the config is updated.
	//
	// The function is called with the new config settings.
	OnConfigUpdated(settings T)
}

type Config[T any] interface {
	SetDefaults(settings T)
	Load() (T, error)
	Save(settings T) error
	EnsureConfig() error
	ConfigFilePath() string
	ConfigListeners() []ConfigListener[T]
	Subscribe(listener ConfigListener[T])
	Unsubscribe(listener ConfigListener[T])
	SettingsAndAlready() (T, bool)
}

type stdConfig[T any] struct {
	filename    string
	already     bool
	settings    T
	rw          sync.RWMutex
	locker      sync.Mutex
	listeners   []ConfigListener[T]
	listenersRW sync.RWMutex

	viper     *viper.Viper
	isPtrType bool
}

func NewConfig[T any](filename string) Config[T] {

	cfg := &stdConfig[T]{}

	t := reflect.TypeFor[T]()
	cfg.isPtrType = t.Kind() == reflect.Ptr
	cfg.viper = viper.New()
	cfg.filename = filename

	return cfg
}

func (cfg *stdConfig[T]) SetDefaults(settings T) {
	cfg.rw.Lock()
	defer cfg.rw.Unlock()
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
		cfg.viper.SetDefault(field.Name, fv.Interface())
	}
}

func (cfg *stdConfig[T]) Load() (T, error) {
	cfg.locker.Lock()
	defer cfg.locker.Unlock()

	var settings T
	err := cfg.viper.ReadInConfig()
	if _, ok := err.(*fs.PathError); ok || os.IsNotExist(err) {
		err = nil
	}
	if err != nil {
		return settings, err
	}

	if cfg.isPtrType {
		settings = reflect.New(reflect.TypeOf(settings).Elem()).Interface().(T)
		err = cfg.viper.Unmarshal(settings)
	} else {
		err = cfg.viper.Unmarshal(&settings)
	}

	if err == nil {
		onSettingsChanged(cfg, settings)
	}

	return settings, err
}

func (cfg *stdConfig[T]) Save(settings T) error {

	cfg.locker.Lock()
	defer cfg.locker.Unlock()

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
		value := cfg.viper.Get(field.Name)
		if reflect.DeepEqual(value, fv) {
			continue
		}
		cfg.viper.Set(field.Name, fv)
		needSave = true
	}
	if !needSave {
		return nil
	}

	err := cfg.viper.WriteConfig()

	if err == nil {
		onSettingsChanged(cfg, settings)
	}

	return err
}

// EnsureConfig ensures that the configuration directory exists.
// It creates the directory if it does not exist.
// It returns an error if any error occurs during the creation process.
func (cfg *stdConfig[T]) EnsureConfig() error {

	configFilePath := filepath.Join(RootPath(), cfg.filename)
	cfg.viper.SetConfigFile(configFilePath)

	configDirPath := filepath.Dir(configFilePath)
	_, err := os.Stat(configDirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(configDirPath, 0755)
	}

	if err == nil {
		_, err = cfg.Load()
	}

	return err
}

// ConfigFilePath returns the path to the configuration file.
//
// The path is determined by the configuration name passed to NewConfig and
// the root path of the application, which can be set using the environment
func (cfg *stdConfig[T]) ConfigFilePath() string {
	return cfg.viper.ConfigFileUsed()
}

func (cfg *stdConfig[T]) SettingsAndAlready() (T, bool) {
	cfg.rw.RLock()
	defer cfg.rw.RUnlock()
	return cfg.settings, cfg.already
}

func (cfg *stdConfig[T]) ConfigListeners() []ConfigListener[T] {
	cfg.listenersRW.RLock()
	defer cfg.listenersRW.RUnlock()

	return cfg.listeners
}

func (cfg *stdConfig[T]) Subscribe(listener ConfigListener[T]) {
	cfg.listenersRW.Lock()
	defer cfg.listenersRW.Unlock()

	settings, already := cfg.SettingsAndAlready()
	if already {
		listener.OnConfigUpdated(settings)
	}
	cfg.listeners = append(cfg.listeners, listener)
}

func (cfg *stdConfig[T]) Unsubscribe(listener ConfigListener[T]) {
	cfg.listenersRW.Lock()
	defer cfg.listenersRW.Unlock()
	for i, l := range cfg.listeners {
		if l == listener {
			cfg.listeners = append(cfg.listeners[:i], cfg.listeners[i+1:]...)
			break
		}
	}
}

func onSettingsChanged[T any](cfg *stdConfig[T], settings T) {
	cfg.rw.Lock()
	cfg.already = true
	cfg.settings = settings
	cfg.rw.Unlock()

	listeners := cfg.ConfigListeners()
	if len(listeners) > 0 {
		for _, l := range listeners {
			l.OnConfigUpdated(settings)
		}
	}
}
