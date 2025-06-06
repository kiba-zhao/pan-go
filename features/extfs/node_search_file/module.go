// Define runtime module for node search file
package nodesearchfile

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/config"
	"pan/lib/feature"
	"pan/lib/injection"
	"pan/lib/log"
	"pan/lib/repository"
	"pan/lib/runtime"
	"sync"
)

func New(featureName string, provider injection.ComponentStoreProvider) interface{} {
	module := &stdModule{}
	module.ComponentStoreProvider = provider

	cleaner := &stdTaskCleaner{}
	module.cleaner = cleaner

	worker := &stdTaskWorker{}
	module.worker = worker
	cleaner.worker = worker
	worker.reloadChan = make(chan struct{}, 1)

	logger := log.Default()
	worker.logger = logger
	cleaner.logger = logger
	return runtime.NewModule(feature.New(featureName, module), module)
}

type NodeSearchFileConfig = config.Config[*NodeSearchFileSettings]

type stdModule struct {
	injection.ComponentStoreProvider

	Config NodeSearchFileConfig

	worker  *stdTaskWorker
	cleaner *stdTaskCleaner

	controllers     []feature.WebController
	controllersOnce sync.Once

	metaList     []feature.RepositoryMeta
	metaListOnce sync.Once
}

var _ = (feature.WebControllerProvider)((*stdModule)(nil))

func (module *stdModule) WebControllers() []feature.WebController {
	module.controllersOnce.Do(func() {
		module.controllers = []feature.WebController{
			&NodeSearchFileController{},
		}
	})
	return module.controllers
}

var _ = (feature.RepositoryMetaProvider)((*stdModule)(nil))

func (module *stdModule) RepositoryMetaList() []feature.RepositoryMeta {
	module.metaListOnce.Do(func() {
		module.metaList = []feature.RepositoryMeta{
			feature.NewRepositoryMeta[NodeSearchFileRepository](NewNodeSearchFileRepository()),
			feature.NewRepositoryMeta[NodeSearchTaskRepository](NewNodeSearchTaskRepository()),
		}
	})

	return module.metaList
}

var _ = (config.ConfigListener[*NodeSearchFileSettings])((*stdModule)(nil))

func (module *stdModule) OnConfigUpdated(settings *NodeSearchFileSettings) {
	module.worker.SetParallelThreshold(settings.ParallelThreshold)
	module.cleaner.SetLifecycle(settings.Lifecycle)
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (module *stdModule) Components() []injection.Component {

	components := []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
		injection.NewComponent(module.cleaner, injection.ComponentNoneScope),
	}

	components = feature.AppendExternalComponent[TaskWorker](components, module.worker)

	// services
	components = feature.AppendInternalComponent[NodeSearchFileInternalService](components, &NodeSearchFileService{})
	return components
}

var _ = (runtime.ProviderModule)((*stdModule)(nil))

func (module *stdModule) Modules() []interface{} {
	return []interface{}{
		config.NewWithDefaults("extfs_searchfile.toml", newDefaultSettings()),
	}
}

var _ = (repository.Repository)((*stdModule)(nil))

func (module *stdModule) SetupToRepository(db repository.RepositoryDB) error {
	return db.AutoMigrate(&NodeSearchTask{}, &NodeSearchFile{})
}

var _ = (feature.RepositoryDBModule)((*stdModule)(nil))

func (module *stdModule) IsTempDB() bool {
	return true
}

var _ = (bootstrap.ReadyModule)((*stdModule)(nil))

func (module *stdModule) Ready(ctx context.Context) error {
	fileRater := &stdFileRater{}
	module.worker.RegisterFileRater(fileRater)
	defer module.worker.UnregisterFileRater(fileRater)

	module.Config.Subscribe(module)
	defer module.Config.Unsubscribe(module)

	var wg sync.WaitGroup
	wg.Add(2)

	go func(worker *stdTaskWorker) {
		defer wg.Done()
		worker.Run(ctx)
	}(module.worker)

	go func(cleaner *stdTaskCleaner) {
		defer wg.Done()
		cleaner.Run(ctx)
	}(module.cleaner)

	wg.Wait()
	return nil
}
