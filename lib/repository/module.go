package repository

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/injection"
)

func New() interface{} {
	module := &stdRepositoryModule{}

	cluster := &stdRepositoryCluster{}
	module.cluster = cluster

	return module
}

type stdRepositoryModule struct {
	cluster *stdRepositoryCluster
}

var _ = (injection.ComponentProvider)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentNoneScope),
		injection.NewComponent[RepositoryCluster](m.cluster, injection.ComponentExternalScope),
	}
}

var _ = (bootstrap.DeferModule)((*stdRepositoryModule)(nil))

func (m *stdRepositoryModule) Defer(ctx context.Context) error {
	err := m.cluster.InitBaseDB(DBPath())
	if err == nil {
		err = m.cluster.InitTempDB(TempDBPath())
	}
	return err
}
