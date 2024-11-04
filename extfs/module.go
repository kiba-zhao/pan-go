package extfs

import (
	"pan/app"
	"pan/app/bootstrap"
	"pan/app/node"

	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/repositories"
	repoImpl "pan/extfs/repositories/impl"
	"pan/extfs/services"
	"pan/extfs/vfs"
	"sync"
)

func New() interface{} {
	return app.NewSample(&module{vfs: &vfs.VFS{}})
}

const moduleName = "extfs"

type module struct {
	Node        node.NodeModule
	DBProvider  app.RepositoryDBProvider
	NodeScope   node.NodeScopeModule
	controllers []interface{}
	once        sync.Once
	vfs         *vfs.VFS
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
			&controllers.RemoteFileBlockController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&models.NodeItem{},
	}
}

func (m *module) Components() []bootstrap.Component {

	// base
	components := []bootstrap.Component{
		bootstrap.NewComponent(m.DBProvider, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent(m.NodeScope, bootstrap.ComponentInternalScope),
	}

	// services
	components = app.AppendSampleInternalComponent[services.NodeItemInternalService](components, &services.NodeItemService{})
	components = app.AppendSampleInternalComponent[services.FileItemInternalService](components, &services.FileItemService{})

	components = app.AppendSampleComponent(components, &services.RemoteNodeItemService{})
	components = app.AppendSampleComponent(components, &services.RemoteNodeService{Provider: m})
	components = app.AppendSampleComponent(components, &services.RemoteFileItemService{})
	components = app.AppendSampleComponent(components, &services.RemoteFileBlockService{})

	// repositories
	components = app.AppendSampleComponent[repositories.NodeItemRepository](components, &repoImpl.NodeItemRepository{})

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, bootstrap.NewComponent(ctrl, bootstrap.ComponentNoneScope))
	}

	// vfs components
	vfsComponents := m.vfs.VFSComponents()
	components = append(components, vfsComponents...)
	return components
}

func (m *module) NodeManager() node.NodeManager {
	if m.Node == nil {
		return nil
	}
	mgr := m.Node.NodeManager()
	return mgr
}

func (m *module) Modules() []interface{} {
	return []interface{}{
		m.vfs,
	}
}
