package nodesearchfile

import (
	"os"
	appConfig "pan/app/config"
	"pan/app/injection"
	"pan/app/sample"
	"pan/runtime"
	"path"
	"reflect"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Agent interface {
	Settings() NodeSearchFileSettings
	FileRaters() []FileRater
}

type agentImpl struct {
	Config appConfig.Config[*NodeSearchFileSettings]

	settings   NodeSearchFileSettings
	settingsRW sync.RWMutex

	registry   runtime.Registry
	registryRW sync.RWMutex

	db     *gorm.DB
	dbOnce sync.Once

	provider injection.ComponentStoreProvider
}

// FileRaters returns a slice of all file raters registered in the component store.
//
// It returns an empty slice if the component store is not initialized yet.
func (agent *agentImpl) FileRaters() []FileRater {
	agent.registryRW.RLock()
	registry := agent.registry
	agent.registryRW.RUnlock()

	if registry == nil {
		return nil
	}

	raters := runtime.ModulesForType[FileRater](registry)
	return raters
}

// Settings returns the current settings of the node search file agent.
//
// It returns a copy of the current settings. If the agent's configuration is
// updated, the returned settings will not be updated.
func (agent *agentImpl) Settings() NodeSearchFileSettings {
	agent.settingsRW.RLock()
	defer agent.settingsRW.RUnlock()
	return agent.settings
}

// OnConfigUpdated updates the agent's settings with the given settings.
//
// It updates the agent's local copy of the settings and does not return an error.
func (agent *agentImpl) OnConfigUpdated(settings *NodeSearchFileSettings) {
	agent.settingsRW.Lock()
	defer agent.settingsRW.Unlock()
	agent.settings = *settings
}

func (agent *agentImpl) Init(registry runtime.Registry) error {

	agent.registryRW.Lock()
	agent.registry = registry
	agent.registryRW.Unlock()

	return nil
}

func (agent *agentImpl) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[FileRater](),
	}
}

// DB returns the Gorm DB instance used by the agent.
//
// It initializes the database by attempting to create the database file in the
// path specified by the DBPath setting. If the path does not exist, it creates
// the directory recursively with 0755 permissions. If the database file already
// exists, it is used without modification. If the database file does not exist,
// it creates a new database using the given path and sets up the necessary
// tables using AutoMigrate. It returns the initialized Gorm DB instance.
func (agent *agentImpl) DB() *gorm.DB {
	agent.dbOnce.Do(func() {
		settings, err := agent.Config.Read()
		if err != nil {
			panic(err)
		}
		configPath := settings.DBPath
		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(configPath, 0755)
		}
		if err == nil {
			dbPath := path.Join(configPath, "extfs-filerate.db")
			agent.db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		}
		if err == nil {
			_ = agent.db.AutoMigrate(&NodeSearchTask{})
		}
	})
	return agent.db
}

// ComponentStore returns the component store that the agent uses to
// store its components. This is the same component store that is
// passed to the agent's constructor.
func (agent *agentImpl) ComponentStore() injection.ComponentStore {
	return agent.provider.ComponentStore()
}

func (agent *agentImpl) Components() []injection.Component {

	components := []injection.Component{}

	// repositories
	components = sample.AppendSampleComponent(components, NewNodeSearchTaskRepository(agent.DB()))
	components = sample.AppendSampleComponent(components, NewNodeSearchFileRepository(agent.DB()))

	return components
}
