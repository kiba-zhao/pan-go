package config

import (
	"pan/pkg/log"
	"slices"
	"sync"
)

type ConfigurerListener[T any] interface {
	OnConfigUpdated(config T)
}

type Configurer[T any] interface {
	Configure(config T) error
	Config() T
	ConfigListeners() []ConfigurerListener[T]
	Subscribe(listener ConfigurerListener[T])
	Unsubscribe(listener ConfigurerListener[T])
}

type stdConfigurer[T any] struct {
	logger log.Logger

	config   T
	configRW sync.RWMutex

	listeners   []ConfigurerListener[T]
	listenersRW sync.RWMutex
}

func NewConfigurer[T any](logger log.Logger) Configurer[T] {
	return &stdConfigurer[T]{
		logger: logger,
	}
}

func NewConfigurerWithConfig[T any](logger log.Logger, config T) Configurer[T] {
	configurer := NewConfigurer[T](logger)
	configurer.Configure(config)
	return configurer
}

func (configurer *stdConfigurer[T]) Configure(config T) error {
	configurer.configRW.Lock()
	configurer.config = config
	configurer.configRW.Unlock()

	listeners := configurer.ConfigListeners()
	if len(listeners) > 0 {
		for _, ln := range listeners {
			ln.OnConfigUpdated(config)
		}
	}

	return nil
}

func (configurer *stdConfigurer[T]) Config() T {
	configurer.configRW.RLock()
	defer configurer.configRW.RUnlock()

	return configurer.config
}

func (configurer *stdConfigurer[T]) ConfigListeners() []ConfigurerListener[T] {
	configurer.listenersRW.RLock()
	defer configurer.listenersRW.RUnlock()

	return slices.Clone(configurer.listeners)
}

func (configurer *stdConfigurer[T]) Subscribe(listener ConfigurerListener[T]) {
	configurer.listenersRW.Lock()
	configurer.listeners = append(configurer.listeners, listener)
	configurer.listenersRW.Unlock()

	listener.OnConfigUpdated(configurer.Config())

}

func (configurer *stdConfigurer[T]) Unsubscribe(listener ConfigurerListener[T]) {
	configurer.listenersRW.Lock()
	defer configurer.listenersRW.Unlock()
	for i, l := range configurer.listeners {
		if l == listener {
			configurer.listeners = append(configurer.listeners[:i], configurer.listeners[i+1:]...)
			break
		}
	}
}
