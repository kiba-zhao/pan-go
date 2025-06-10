package config

import (
	"pan/lib/injection"
)

type stdConfigModule[T any] struct {
	cfg Config[T]
}

func New[T any](filename string) interface{} {
	cfgModule := &stdConfigModule[T]{}
	cfgModule.cfg = NewConfig[T](filename)
	return cfgModule
}

func NewWithDefaults[T any](filename string, settings T) interface{} {
	cfgModule := &stdConfigModule[T]{}
	cfgModule.cfg = NewConfig[T](filename)
	cfgModule.cfg.SetDefaults(settings)
	return cfgModule
}

func (m *stdConfigModule[T]) Defer() error {
	return m.cfg.EnsureConfig()
}

func (m *stdConfigModule[T]) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m.cfg, injection.ComponentExternalScope),
	}
}
