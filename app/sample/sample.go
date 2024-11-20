package sample

import (
	"os"
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/peer"
	"pan/app/web"
	"path"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

func New[T SampleProvider](provider T) interface{} {
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

func (s *sample[T]) PeerScope() []byte {
	return []byte(s.provider.Name() + ".")
}

func (s *sample[T]) PeerAppModules() []peer.PeerAppModule {
	modules := make([]peer.PeerAppModule, 0)
	for _, c := range s.provider.Controllers() {
		if m, ok := c.(peer.PeerAppModule); ok {
			modules = append(modules, m)
		}
	}
	return modules
}

func (s *sample[T]) WebScope() string {
	return "/api/" + s.provider.Name()
}

func (s *sample[T]) WebControllers() []web.WebController {
	controllers := make([]web.WebController, 0)
	for _, c := range s.provider.Controllers() {
		if m, ok := c.(web.WebController); ok {
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
		bootstrap.NewComponent[peer.PeerScopeModule](s, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent[web.WebScopeModule](s, bootstrap.ComponentInternalScope),
		bootstrap.NewComponent(s.provider, bootstrap.ComponentNoneScope),
	}
}

func (s *sample[T]) Modules() []interface{} {
	return []interface{}{s.provider}
}
