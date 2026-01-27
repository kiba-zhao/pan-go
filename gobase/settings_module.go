package gobase

import (
	"context"
	"pan/internal/bootstrap"
	"pan/internal/injection"
	"pan/internal/settings"
)

func NewSettingsModule(cfg settings.SettingsConfig) interface{} {
	module := &stdSettingsProxyModule{}
	module.cfg = cfg
	return module
}

type stdSettingsProxyModule struct {
	Configurer settings.SettingsConfigurer
	cfg        settings.SettingsConfig
}

var _ = (injection.ComponentProvider)((*stdSettingsProxyModule)(nil))

func (module *stdSettingsProxyModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
	}
}

var _ = (bootstrap.DeferModule)((*stdSettingsProxyModule)(nil))

func (module *stdSettingsProxyModule) Defer(ctx context.Context) error {
	var err error
	if module.Configurer != nil {
		err = module.Configurer.Configure(module.cfg)
	}
	return err
}
