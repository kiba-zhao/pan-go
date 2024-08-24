package app

import (
	"os"
	"pan/app/config"
	"pan/app/net"
	"pan/app/node"
	"pan/app/repositories"
	"pan/runtime"
	"path"
	"reflect"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type RepositoryDB = repositories.RepositoryDB

type RepositoryDBProvider = repositories.RepositoryDBProvider

var DBForProvider = repositories.DBForProvider

type SampleProvider interface {
	Name() string
	Controllers() []interface{}
	Models() []interface{}
}

type sample[T SampleProvider] struct {
	provider T
	rootPath string
	db       RepositoryDB
	dbLocker sync.RWMutex
}

func NewSample[T SampleProvider](provider T) interface{} {
	return &sample[T]{provider: provider}
}

func (s *sample[T]) DB() RepositoryDB {
	s.dbLocker.RLock()
	defer s.dbLocker.RUnlock()
	return s.db
}

func (s *sample[T]) OnConfigUpdated(settings config.AppSettings) {

	s.dbLocker.Lock()
	defer s.dbLocker.Unlock()

	if s.rootPath == settings.RootPath {
		return
	}

	s.rootPath = settings.RootPath

	models := s.provider.Models()
	if len(models) <= 0 {
		return
	}

	var db RepositoryDB
	_, err := os.Stat(s.rootPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(s.rootPath, 0755)
	}
	if err == nil {
		dbPath := path.Join(s.rootPath, s.provider.Name()+".db")
		db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
	if err == nil {
		_ = db.AutoMigrate(models...)
	}
	s.db = db
}

func (s *sample[T]) NodeProviderName() []byte {
	return []byte(s.provider.Name() + ".")
}

func (s *sample[T]) NodeAppModule() []node.NodeAppModule {
	modules := make([]node.NodeAppModule, 0)
	for _, c := range s.provider.Controllers() {
		if m, ok := c.(node.NodeAppModule); ok {
			modules = append(modules, m)
		}
	}
	return modules
}

func (s *sample[T]) WebProviderName() string {
	return "/api/" + s.provider.Name()
}

func (s *sample[T]) WebControllers() []net.WebController {
	controllers := make([]net.WebController, 0)
	for _, c := range s.provider.Controllers() {
		if m, ok := c.(net.WebController); ok {
			controllers = append(controllers, m)
		}
	}
	return controllers
}

func (s *sample[T]) Models() []interface{} {
	return []interface{}{
		s.provider,
	}
}

func (s *sample[T]) Components() []runtime.Component {
	return []runtime.Component{
		runtime.NewComponent[RepositoryDBProvider](s, runtime.ComponentInternalScope),
		runtime.NewComponent(s.provider, runtime.ComponentNoneScope),
	}
}

func (s *sample[T]) Modules() []interface{} {
	return []interface{}{s.provider}
}

func AppendSampleComponent[T any](components []runtime.Component, component T) []runtime.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, runtime.NewComponentByType(reflect.TypeOf(component), component, runtime.ComponentNoneScope))
	}
	components = append(components, runtime.NewComponent(component, runtime.ComponentInternalScope))
	return components
}
