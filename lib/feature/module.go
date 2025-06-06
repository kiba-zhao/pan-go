package feature

import (
	"pan/lib/bootstrap"
	"pan/lib/injection"
	"pan/lib/repository"
)

type stdFeatureModule struct {
	RepositoryCluster repository.RepositoryCluster

	featureHelper *stdFeatureHelper
	brokerHelper  *stdBrokerHelper
}

func New(name string, feature interface{}) interface{} {
	module := &stdFeatureModule{}

	brokerHelper := &stdBrokerHelper{}
	module.brokerHelper = brokerHelper

	helper := &stdFeatureHelper{}
	module.featureHelper = helper
	brokerHelper.featureHelper = helper
	helper.feature = feature
	helper.name = name

	return module
}

var _ = (injection.ComponentProvider)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) Components() []injection.Component {
	brokerHelper := module.brokerHelper
	components := []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
		injection.NewComponent[BrokerHelper](brokerHelper, injection.ComponentInternalScope),
	}

	featureHelper := module.featureHelper
	controllers := getWebControllers(featureHelper.feature)
	if len(controllers) > 0 {
		for _, ctrl := range controllers {
			components = append(components, injection.NewComponent(ctrl, injection.ComponentNoneScope))
		}
	}

	metaList := getRepositoryMetaList(featureHelper.feature)
	if len(metaList) > 0 {
		for _, meta := range metaList {
			components = append(components, injection.NewComponentByType(meta.Type, meta.Repository, injection.ComponentInternalScope))
		}
	}

	topics := getPeerTopics(featureHelper.feature)
	if len(topics) > 0 {
		for _, topic := range topics {
			components = append(components, injection.NewComponent(topic, injection.ComponentNoneScope))
		}
	}

	return components
}

var _ = (injection.ComponentStoreProvider)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) ComponentStore() injection.ComponentStore {
	featureHelper := module.featureHelper
	if storeProvider, ok := featureHelper.feature.(injection.ComponentStoreProvider); ok {
		return storeProvider.ComponentStore()
	}
	return nil
}

var _ = (bootstrap.DeferModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) Defer() error {
	featureHelper := module.featureHelper
	repository := getRepository(featureHelper.feature)
	if repository == nil {
		return nil
	}

	if isTempDB(featureHelper.feature) {
		return module.RepositoryCluster.AttachTempModule(module)
	}
	return module.RepositoryCluster.AttachBaseModule(module)
}

var _ = (bootstrap.DestroyModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) Destroy() {
	featureHelper := module.featureHelper
	repository := getRepository(featureHelper.feature)
	if repository == nil {
		return
	}

	if isTempDB(featureHelper.feature) {
		module.RepositoryCluster.DetachTempModule(module)
	} else {
		module.RepositoryCluster.DetachBaseModule(module)
	}
}

type RepositoryDBModule interface {
	IsTempDB() bool
}

func isTempDB(feature interface{}) bool {
	dbModule, ok := feature.(RepositoryDBModule)
	if ok {
		return dbModule.IsTempDB()
	}
	return ok
}
