package app

import (
	"os"
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/net"
	"pan/app/node"
	"pan/app/repositories"
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
	Config   config.AppConfig
	provider T
	db       RepositoryDB
	once     sync.Once
}

func NewSample[T SampleProvider](provider T) interface{} {
	return &sample[T]{provider: provider}
}

func (s *sample[T]) DB() RepositoryDB {
	s.once.Do(func() {
		models := s.provider.Models()
		if len(models) <= 0 {
			return
		}
		configPath := path.Dir(s.Config.ConfigFilePath())
		var db RepositoryDB
		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(configPath, 0755)
		}
		if err == nil {
			dbPath := path.Join(configPath, s.provider.Name()+".db")
			db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		}
		if err == nil {
			_ = db.AutoMigrate(models...)
		}
		s.db = db
	})
	return s.db
}

func (s *sample[T]) NodeScope() []byte {
	return []byte(s.provider.Name() + ".")
}

func (s *sample[T]) NodeAppModules() []node.NodeAppModule {
	modules := make([]node.NodeAppModule, 0)
	for _, c := range s.provider.Controllers() {
		if m, ok := c.(node.NodeAppModule); ok {
			modules = append(modules, m)
		}
	}
	return modules
}

func (s *sample[T]) WebScope() string {
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

func (s *sample[T]) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent(s, bootstrap.ComponentNoneScope),
		bootstrap.NewComponent[RepositoryDBProvider](s, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent[node.NodeScopeModule](s, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent[net.WebScopeModule](s, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent(s.provider, bootstrap.ComponentNoneScope),
	}
}

func (s *sample[T]) Modules() []interface{} {
	return []interface{}{s.provider}
}

func AppendSampleComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentNoneScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentInternalScope))
	return components
}

func AppendSampleExternalComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentInternalScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentExternalScope))
	return components
}

func AppendSampleInternalComponent[T any](components []bootstrap.Component, component T) []bootstrap.Component {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Interface {
		components = append(components, bootstrap.NewComponentByType(reflect.TypeOf(component), component, bootstrap.ComponentInternalScope))
	}
	components = append(components, bootstrap.NewComponent(component, bootstrap.ComponentInternalScope))
	return components
}
