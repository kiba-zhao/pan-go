package extfs

import (
	"os"
	"pan/app"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"pan/extfs/services"
	"path"
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
	components  []app.Component
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
			err = db.AutoMigrate(&models.Target{})
		}
		if err != nil {
			panic(err)
		}
		m.db = db
	})
	return m.db
}

func (m *module) Components() []app.Component {

	m.ctrlOnce.Do(func() {

		// base
		m.components = append(m.components,
			app.NewComponent(m, app.ComponentNoneScope),
		)

		// repositories
		setupRepository(m, repositories.NewTargetRepository)

		// services
		m.components = append(m.components,
			app.NewComponent(&services.TargetService{}, app.ComponentExternalScope),
			app.NewComponent(&services.DiskFileService{}, app.ComponentExternalScope),
		)

		// controllers
		setupController(m, &controllers.TargetController{})
		setupController(m, &controllers.DiskFileController{})

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
	m.components = append(m.components, app.NewComponent(controller, app.ComponentNoneScope))
}

type LazyRepositoryFunc[T any] func(db *gorm.DB) T

func setupRepository[T any](m *module, repoFunc LazyRepositoryFunc[T]) {
	lazyFunc := func() T { return repoFunc(m.DB()) }
	m.components = append(m.components, app.NewLazyComponent(lazyFunc, app.ComponentInternalScope))
}
