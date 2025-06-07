package repository

import (
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/log"
)

func New() interface{} {
	module := &stdRepositoryModule{}

	cluster := &stdRepositoryCluster{}
	module.cluster = cluster

	return module
}

type stdRepositoryModule struct {
	AppConfig config.AppConfig

	logger  log.Logger
	cluster *stdRepositoryCluster
}

var _ = (injection.ComponentProvider)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
		injection.NewComponent[RepositoryCluster](m.cluster, injection.ComponentExternalScope),
	}
}

var _ = (config.AppConfigListener)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) OnConfigUpdated(settings config.AppSettings) {
	err := m.cluster.InitBaseDB(settings.DBPath)
	if err != nil {
		m.logger.Error("RepositoryModule", "InitBaseDB Error: "+err.Error())
	}

	err = m.cluster.InitTempDB(settings.TempPath)
	if err != nil {
		m.logger.Error("RepositoryModule", "InitTempDB Error: "+err.Error())
	}
}

var _ = (bootstrap.DeferModule)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Defer() error {
	m.AppConfig.Subscribe(m)
	return nil
}

var _ = (bootstrap.DestroyModule)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Destroy() {
	m.AppConfig.Unsubscribe(m)
}
