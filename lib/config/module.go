package config

import (
	"pan/lib/injection"
	"path/filepath"
)

type stdConfigModule[T any] struct {
	cfg  Config[T]
	name string
}

func New[T any](name string) interface{} {
	cfgModule := &stdConfigModule[T]{}
	cfgModule.name = name
	cfgModule.cfg = NewConfig[T]()
	return cfgModule
}

func NewWithDefaults[T any](name string, settings T) interface{} {
	cfgModule := &stdConfigModule[T]{}
	cfgModule.name = name
	cfgModule.cfg = NewConfig[T]()
	cfgModule.cfg.SetDefaults(settings)
	return cfgModule
}

func (m *stdConfigModule[T]) Defer() error {
	return m.cfg.EnsureConfig(filepath.Join(getRootPath(), m.name))
}

func (m *stdConfigModule[T]) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m.cfg, injection.ComponentExternalScope),
	}
}
