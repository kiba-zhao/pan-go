package extfs

import (
	"os"
	"pan/app"
	"pan/extfs/controllers"
	"pan/extfs/dispatchers"
	dispatcherImpl "pan/extfs/dispatchers/impl"
	"pan/extfs/models"
	"pan/extfs/repositories"
	repoImpl "pan/extfs/repositories/impl"
	"pan/extfs/services"
	"pan/runtime"
	"path"
	"reflect"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type webController interface {
	Init(router app.WebRouter)
}

type module struct {
	Config      app.AppConfig
	db          *gorm.DB
	dbOnce      sync.Once
	controllers []webController
	ctrlOnce    sync.Once
}

func New() interface{} {
	return &module{}
}

func (m *module) DB() *gorm.DB {
	m.dbOnce.Do(func() {
		var db *gorm.DB
		settings, err := m.Config.Read()
		if err == nil {
			basePath := settings.RootPath
			_, err := os.Stat(basePath)
			if os.IsNotExist(err) {
				err = os.MkdirAll(basePath, 0755)
			}
			if err == nil {
				dbPath := path.Join(basePath, "extfs.db")
				db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			}
		}
		if err == nil {
			err = db.AutoMigrate(&models.Target{}, &models.TargetFile{})
		}
		if err != nil {
			panic(err)
		}
		m.db = db
	})
	return m.db
}

func (m *module) Controllers() []webController {
	m.ctrlOnce.Do(func() {
		m.controllers = append(m.controllers,
			&controllers.TargetController{},
			&controllers.DiskFileController{},
			&controllers.TargetFileController{},
		)
	})
	return m.controllers
}

func (m *module) Components() []runtime.Component {

	var components []runtime.Component
	// base
	components = append(components,
		runtime.NewComponent(m, runtime.ComponentNoneScope),
		runtime.NewLazyComponent(m.DB, runtime.ComponentInternalScope),
	)

	// repositories
	components = setupComponent[repositories.TargetRepository](components, &repoImpl.TargetRepository{})
	components = setupComponent[repositories.TargetFileRepository](components, &repoImpl.TargetFileRepository{})

	// services
	components = setupComponent(components, &services.TargetService{})
	components = setupComponent(components, &services.DiskFileService{})
	components = setupComponent(components, &services.TargetFileService{})

	// dispatchers
	components = setupComponent(components, dispatcherImpl.NewTargetDispatcherBucket())
	components = setupComponent[dispatchers.TargetDispatcher](components, &dispatcherImpl.TargetDispatcher{})

	// controllers
	for _, ctrl := range m.Controllers() {
		components = append(components, runtime.NewComponent(ctrl, runtime.ComponentNoneScope))
	}

	return components
}

func (m *module) SetupToWeb(webApp app.WebApp) error {

	router := webApp.Group("/api/extfs")

	for _, ctrl := range m.Controllers() {
		ctrl.Init(router)
	}

	return nil
}

func setupComponent[T any](components []runtime.Component, component T) []runtime.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, runtime.NewComponentByType(reflect.TypeOf(component), component, runtime.ComponentNoneScope))
	}
	components = append(components, runtime.NewComponent(component, runtime.ComponentInternalScope))
	return components
}
