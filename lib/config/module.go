package config

import (
	"pan/lib/injection"
	"pan/lib/log"
)

type stdConfigModule[T any] struct {
	configurer Configurer[T]
}

func New[T any]() interface{} {
	logger := log.Default()

	cfgModule := &stdConfigModule[T]{}
	cfgModule.configurer = NewConfigurer[T](logger)
	return cfgModule
}

func (m *stdConfigModule[T]) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
	}
}
