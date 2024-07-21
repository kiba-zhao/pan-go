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
	rootPath    string
	db          repositories.RepositoryDB
	dbLocker    sync.RWMutex
	controllers []webController
	ctrlOnce    sync.Once
}

func New() interface{} {
	return &module{}
}

func (m *module) OnConfigUpdated(settings app.AppSettings) {

	m.dbLocker.Lock()
	defer m.dbLocker.Unlock()

	if m.rootPath == settings.RootPath {
		return
	}

	m.rootPath = settings.RootPath

	var db repositories.RepositoryDB
	_, err := os.Stat(m.rootPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(m.rootPath, 0755)
	}
	if err == nil {
		dbPath := path.Join(m.rootPath, "extfs.db")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
	if err == nil {
		err = db.AutoMigrate(&models.Target{}, &models.TargetFile{})
	}
	m.db = db
}

func (m *module) DB() repositories.RepositoryDB {
	m.dbLocker.RLock()
	defer m.dbLocker.RUnlock()
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
		runtime.NewComponent[repositories.ComponentProvider](m, runtime.ComponentInternalScope),
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
