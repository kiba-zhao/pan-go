package extfs

import (
	"os"
	"pan/app"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"pan/extfs/services"
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

func (m *module) GetComponents() []interface{} {

	m.ctrlOnce.Do(func() {
		m.controllers = []webController{
			&controllers.TargetController{},
			&controllers.DiskFileController{},
		}
	})

	components := []interface{}{m}
	for _, ctrl := range m.controllers {
		components = append(components, ctrl)
	}
	return components
}

func (m *module) GetDependencies() map[reflect.Type]interface{} {

	m.dbOnce.Do(func() {
		settings, err := m.Config.Read()
		if err == nil {
			basePath := settings.RootPath
			_, err := os.Stat(basePath)
			if os.IsNotExist(err) {
				err = os.MkdirAll(basePath, 0755)
			}
			if err == nil {
				dbPath := path.Join(basePath, "extfs.db")
				m.db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
			}
		}
		if err == nil {
			err = m.db.AutoMigrate(&models.Target{})
		}
		if err != nil {
			panic(err)
		}
	})

	return map[reflect.Type]interface{}{
		reflect.TypeFor[repositories.TargetRepository](): repositories.NewTargetRepository(m.db),
		reflect.TypeFor[*services.TargetService]():       &services.TargetService{},
		reflect.TypeFor[*services.DiskFileService]():     &services.DiskFileService{},
	}
}

func (m *module) SetupToWeb(webApp app.WebApp) error {

	router := webApp.Group("/api/extfs")

	for _, ctrl := range m.controllers {
		ctrl.Init(router)
	}

	return nil
}
