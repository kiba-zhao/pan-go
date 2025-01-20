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

func (agent *agentImpl) Settings() NodeSearchFileSettings {
	agent.settingsRW.RLock()
	defer agent.settingsRW.RUnlock()
	return agent.settings
}

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
