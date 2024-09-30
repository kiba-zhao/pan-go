package extfs

import (
	"pan/app"
	"pan/app/node"

	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/repositories"
	repoImpl "pan/extfs/repositories/impl"
	"pan/extfs/services"
	"pan/runtime"
	"sync"
)

func New() interface{} {
	return app.NewSample(&module{})
}

const moduleName = "extfs"

type module struct {
	Node        node.NodeModule
	DBProvider  app.RepositoryDBProvider
	controllers []interface{}
	once        sync.Once
}

func (m *module) Name() string {
	return moduleName
}

func (m *module) Controllers() []interface{} {
	m.once.Do(func() {

		m.controllers = []interface{}{
			&controllers.NodeItemController{},
			&controllers.RemoteNodeController{},
			&controllers.FileItemController{},
			&controllers.RemoteNodeItemController{},
			&controllers.RemoteFileItemController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&models.NodeItem{},
	}
}

func (m *module) Components() []runtime.Component {

	// base
	components := []runtime.Component{
		runtime.NewComponent(m.DBProvider, runtime.ComponentInternalScope),
	}

	// services
	components = app.AppendSampleInternalComponent[services.NodeItemInternalService](components, &services.NodeItemService{})
	components = app.AppendSampleComponent(components, &services.RemoteNodeItemService{})
	components = app.AppendSampleComponent(components, &services.RemoteNodeService{Provider: m})
	components = app.AppendSampleComponent(components, &services.FileItemService{})
	components = app.AppendSampleComponent(components, &services.RemoteFileItemService{})

	// repositories
	components = app.AppendSampleComponent[repositories.NodeItemRepository](components, &repoImpl.NodeItemRepository{})

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, runtime.NewComponent(ctrl, runtime.ComponentNoneScope))
	}

	return components
}

func (m *module) NodeManager() node.NodeManager {
	if m.Node == nil {
		return nil
	}
	mgr := m.Node.NodeManager()
	return mgr
}
