// Define utils for building simple application modules
package sample

import (
	"os"
	"pan/lib/config"
	"pan/lib/injection"
	"pan/lib/peer"
	"pan/lib/web"
	"path/filepath"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Define sample provider
//
// It used to provide sample module name and models
type SampleProvider interface {
	// Name returns sample module name
	Name() string
	// Models returns sample module models
	Models() []interface{}
}

// Define sample interface
type Sample interface {
	// DB returns sample module database instance
	DB() RepositoryDB
	// extends SamplePeer interface
	SamplePeer
}

type sample[T SampleProvider] struct {
	Config     config.AppConfig
	PeerModule peer.PeerModule
	provider   T
	db         RepositoryDB
	once       sync.Once
}

// New returns a new sample module with the given provider.
//
// It is used to create a new sample module, which will have the provider
// as its provider.
func New[T SampleProvider](provider T) Sample {
	return &sample[T]{provider: provider}
}

// DB returns a database instance associated with the sample module.
//
// It automatically creates the database file in the same directory as the
// configuration file if it does not exist. It then auto-migrates the models
// provided by the sample module provider. If the provider does not provide any
// models, it simply returns a nil database instance.
// It implements the Sample interface.
func (s *sample[T]) DB() RepositoryDB {
	s.once.Do(func() {
		models := s.provider.Models()
		if len(models) <= 0 {
			return
		}
		configPath := filepath.Dir(s.Config.ConfigFilePath())
		var db RepositoryDB
		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(configPath, 0755)
		}
		if err == nil {
			dbPath := filepath.Join(configPath, s.provider.Name()+".db")
			db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		}
		if err == nil {
			_ = db.AutoMigrate(models...)
		}
		s.db = db
	})
	return s.db
}

// PeerScope returns a scope for the p2p application associated with the sample module.
//
// It is used to identify the sample module in the p2p application.
// It implements the PeerScopeModule interface.
func (s *sample[T]) PeerScope() []byte {
	return []byte(s.provider.Name() + ".")
}

// PeerAppModules returns a slice of PeerAppModules provided by the sample module's provider.
//
// It checks if the provider implements the PeerAppModuleProvider interface and,
// if so, retrieves the PeerAppModules from the provider. If the provider does
// not implement the interface, it returns nil.
// It implements the PeerAppModuleProvider interface.
func (s *sample[T]) PeerAppModules() []peer.PeerAppModule {
	provider, ok := any(s.provider).(peer.PeerAppModuleProvider)
	if !ok {
		return nil
	}
	return provider.PeerAppModules()
}

// WebScope returns the web scope for the sample module.
//
// The web scope is used to namespace the web controllers provided by the sample module.
// It is used by the web engine to register the controllers.
// It implements the WebScopeModule interface.
func (s *sample[T]) WebScope() string {
	return "/api/" + s.provider.Name()
}

// WebControllers returns a slice of WebControllers provided by the sample module's provider.
//
// It checks if the provider implements the WebControllerProvider interface and,
// if so, retrieves the WebControllers from the provider. If the provider does
// not implement the interface, it returns nil.
// It implements the WebControllerProvider interface.

func (s *sample[T]) WebControllers() []web.WebController {
	provider, ok := any(s.provider).(web.WebControllerProvider)
	if !ok {
		return nil
	}
	return provider.WebControllers()
}

// Components returns a slice of injection.Component that are provided by the
// sample module. This is just the engine itself, and the scope is set to
// ComponentNoneScope as the sample module is not injected into any other
// components.
func (s *sample[T]) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(s, injection.ComponentNoneScope),
	}
}

// Modules returns a slice of sub-modules.
//
// This method is part of the Engine's ProviderModule interface, and is used
// to register the sub-modules in the engine during Mount.
//
// It returns the sample module provider itself as a sub-module.
func (s *sample[T]) Modules() []interface{} {
	return []interface{}{s.provider}
}
