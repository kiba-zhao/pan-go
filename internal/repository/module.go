package repository

import (
	"context"
	"pan/internal/bootstrap"
	"pan/internal/config"
	"pan/internal/injection"
	"pan/internal/log"
)

func New() interface{} {

	manager := &stdRepositoryManager{}

	logger := log.Default()
	manager.logger = logger

	configurer := config.NewConfigurer[RepositoryConfig](logger)

	module := &stdRepositoryModule{}
	module.logger = logger
	module.manager = manager
	module.configurer = configurer

	return module
}

type stdRepositoryModule struct {
	logger     log.Logger
	manager    *stdRepositoryManager
	configurer RepositoryConfigurer
}

var _ = (injection.ComponentProvider)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
		injection.NewComponent[RepositoryManager](m.manager, injection.ComponentExternalScope),
		injection.NewComponent(m.configurer, injection.ComponentExternalScope),
	}
}

var _ = (RepositoryConfigListener)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) OnConfigUpdated(cfg RepositoryConfig) {
	if cfg == nil {
		m.logger.Error("RepositoryModule", "OnConfigUpdated Error: Invalid RepositoryConfig")
		return
	}
	err := m.manager.Setup(cfg)
	if err != nil {
		m.logger.Error("RepositoryModule", "OnConfigUpdated Error: "+err.Error())
	}
}

var _ = (bootstrap.DeferModule)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Defer(ctx context.Context) error {
	m.configurer.Subscribe(m)
	return nil
}

var _ = (bootstrap.DestroyModule)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Destroy() {
	m.configurer.Unsubscribe(m)
}
