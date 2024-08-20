package extfs

import (
	"pan/app"

	"pan/extfs/controllers"
	"pan/extfs/dispatchers"
	dispatcherImpl "pan/extfs/dispatchers/impl"
	"pan/extfs/models"
	"pan/extfs/repositories"
	repoImpl "pan/extfs/repositories/impl"
	"pan/extfs/services"
	"pan/runtime"
	"reflect"
	"sync"
)

func New() interface{} {
	return app.NewSample(&module{})
}

const moduleName = "extfs"

type module struct {
	DBProvider  app.RepositoryDBProvider
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
			&controllers.TargetController{},
			&controllers.TargetFileController{},
		}
	})
	return m.controllers
}

func (m *module) Models() []interface{} {
	return []interface{}{
		&models.Target{}, &models.TargetFile{},
	}
}

func (m *module) Components() []runtime.Component {

	// base
	components := []runtime.Component{
		runtime.NewComponent(m.DBProvider, runtime.ComponentInternalScope),
	}

	// repositories
	components = setupComponent[repositories.TargetRepository](components, &repoImpl.TargetRepository{})
	components = setupComponent[repositories.TargetFileRepository](components, &repoImpl.TargetFileRepository{})

	// services
	components = setupComponent(components, &services.TargetService{})
	components = setupComponent(components, &services.TargetFileService{})

	// dispatchers
	components = setupComponent[dispatchers.TargetDispatcher](components, dispatcherImpl.NewTargetDispatcher())

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, runtime.NewComponent(ctrl, runtime.ComponentNoneScope))
	}

	return components
}

func setupComponent[T any](components []runtime.Component, component T) []runtime.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, runtime.NewComponentByType(reflect.TypeOf(component), component, runtime.ComponentNoneScope))
	}
	components = append(components, runtime.NewComponent(component, runtime.ComponentInternalScope))
	return components
}
