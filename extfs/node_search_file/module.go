package nodesearchfile

import (
	appConfig "pan/app/config"
	"pan/app/injection"
	"pan/app/sample"
)

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

	m.worker = &ratingTaskWorkerImpl{}

	m.fileRater = &fileRaterImpl{}
	return &m
}

type moduleImpl struct {
	agent     Agent
	config    appConfig.Config[*NodeSearchFileSettings]
	worker    NodeSearchTaskWorker
	fileRater FileRater

	provider injection.ComponentStoreProvider
}

func (m *moduleImpl) ComponentStore() injection.ComponentStore {
	return m.provider.ComponentStore()
}

func (m *moduleImpl) Components() []injection.Component {

	components := []injection.Component{
		injection.NewComponent(m.config, injection.ComponentInternalScope),
	}

	components = sample.AppendSampleComponent(components, m.agent)
	components = sample.AppendSampleComponent(components, m.worker)

	return components
}

func (m *moduleImpl) Modules() []interface{} {
	return []interface{}{
		m.config,
		m.agent,
		m.worker,
		m.fileRater,
	}
}
