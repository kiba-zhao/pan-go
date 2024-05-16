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
	initOnce    sync.Once
	components  []runtime.Component
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

func (m *module) Components() []runtime.Component {

	m.initOnce.Do(func() {

		// base
		m.components = append(m.components,
			runtime.NewComponent(m, runtime.ComponentNoneScope),
			runtime.NewLazyComponent(m.DB, runtime.ComponentInternalScope),
		)

		// repositories
		setupComponent[repositories.TargetRepository](m, &repoImpl.TargetRepository{})
		setupComponent[repositories.TargetFileRepository](m, &repoImpl.TargetFileRepository{})

		// services
		setupComponent(m, &services.TargetService{})
		setupComponent(m, &services.DiskFileService{})
		setupComponent(m, &services.TargetFileService{})

		// controllers
		setupController(m, &controllers.TargetController{})
		setupController(m, &controllers.DiskFileController{})
		setupController(m, &controllers.TargetFileController{})

		// dispatchers
		setupComponent(m, dispatcherImpl.NewTargetDispatcherBucket())
		setupComponent[dispatchers.TargetDispatcher](m, &dispatcherImpl.TargetDispatcher{})

	})
	return m.components
}

func (m *module) SetupToWeb(webApp app.WebApp) error {

	router := webApp.Group("/api/extfs")

	for _, ctrl := range m.controllers {
		ctrl.Init(router)
	}

	return nil
}

func setupController[T webController](m *module, controller T) {
	m.controllers = append(m.controllers, controller)
	m.components = append(m.components, runtime.NewComponent(controller, runtime.ComponentNoneScope))
}

func setupComponent[T any](m *module, component T) {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		m.components = append(m.components, runtime.NewComponentByType(reflect.TypeOf(component), component, runtime.ComponentNoneScope))
	}
	m.components = append(m.components, runtime.NewComponent(component, runtime.ComponentInternalScope))
}
