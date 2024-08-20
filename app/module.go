package app

import (
	"pan/app/config"
	"pan/app/controllers"
	"pan/app/net"
	"pan/app/node"
	"pan/app/services"
	"pan/runtime"
	"sync"
)

func New() interface{} {
	return runtime.NewModule(&runtime.Injector{}, config.New(), node.New(), net.New(), NewSample(&module{}))
}

const moduleName = "app"

type module struct {
	DBProvider  RepositoryDBProvider
	controllers []interface{}
	once        sync.Once
}

func (m *module) Name() string {
	return moduleName
}

func (m *module) Controllers() []interface{} {
	m.once.Do(func() {
		// TODO: add web and node controllers
		m.controllers = []interface{}{
			&controllers.NodeController{},
			&controllers.DiskFileController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return nil
}

func (m *module) Components() []runtime.Component {
	// base
	components := []runtime.Component{
		runtime.NewComponent(m.DBProvider, runtime.ComponentInternalScope),
		// services
		runtime.NewComponent(&services.DiskFileService{}, runtime.ComponentInternalScope),
	}
	return components
}
