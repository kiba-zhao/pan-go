// Define runtime module for node search file
package nodesearchfile

import (
	appConfig "pan/app/config"
	"pan/app/injection"
	"pan/app/sample"
)

// New creates a new runtime module for node search file.
//
// It takes a component store provider to create the virtual file system.
// It will load the configuration from "extfs_searchfile.toml" file.
// If the configuration file does not exist, it will panic with an error.
// If the configuration file exists but is invalid, it will panic with an error.
func New(provider injection.ComponentStoreProvider) interface{} {
	var m moduleImpl
	m.provider = provider
	m.agent = &agentImpl{provider: provider}

	config, err := appConfig.NewConfig[*NodeSearchFileSettings]("extfs_searchfile.toml")
	if err != nil {
		panic(err)
	}
	config.SetDefaults(newDefaultSettings())
	m.config = config

	m.worker = &taskWorkerImpl{}
	m.cleaner = &taskCleanerImpl{}

	m.fileRater = &fileRaterImpl{}
	return &m
}

type moduleImpl struct {
	agent     Agent
	config    appConfig.Config[*NodeSearchFileSettings]
	worker    TaskWorker
	cleaner   *taskCleanerImpl
	fileRater FileRater

	provider injection.ComponentStoreProvider
}

// ComponentStore returns the component store that the module uses to
// store its components. This is the same component store that is
// passed to the module's constructor.
func (m *moduleImpl) ComponentStore() injection.ComponentStore {
	return m.provider.ComponentStore()
}

// Components returns a slice of injection.Component representing the components
// provided by the module. It includes the configuration, agent, task worker, task
// cleaner, and file rater. The components are scoped as follows:
//
//   - configuration: internal scope
//   - agent: sample scope
//   - task worker: sample scope
//   - task cleaner: none scope
//   - file rater: none scope
func (m *moduleImpl) Components() []injection.Component {

	components := []injection.Component{
		injection.NewComponent(m.config, injection.ComponentInternalScope),
	}

	components = sample.AppendSampleComponent(components, m.agent)
	components = sample.AppendSampleComponent(components, m.worker)
	components = append(components, injection.NewComponent(m.cleaner, injection.ComponentNoneScope))

	return components
}

// Modules returns a slice of sub-modules.
//
// The sub-modules are the configuration, agent, task worker, task cleaner, and file rater.
func (m *moduleImpl) Modules() []interface{} {
	return []interface{}{
		m.config,
		m.agent,
		m.worker,
		m.cleaner,
		m.fileRater,
	}
}
