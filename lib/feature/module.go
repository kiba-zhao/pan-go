package feature

import (
	"context"
	"pan/lib/bootstrap"
	"pan/lib/injection"
	"pan/lib/repository"
)

type stdFeatureModule struct {
	RepositoryManager repository.RepositoryManager

	featureHelper *stdFeatureHelper
}

func New(name string, feature interface{}) interface{} {
	module := &stdFeatureModule{}

	helper := &stdFeatureHelper{}
	module.featureHelper = helper
	helper.feature = feature
	helper.name = name

	return module
}

var _ = (injection.ComponentProvider)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) Components() []injection.Component {

	components := []injection.Component{
		injection.NewComponent(module, injection.ComponentNoneScope),
		injection.NewComponent(module.featureHelper, injection.ComponentNoneScope),
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
			components = append(components, injection.NewComponentByType(meta.metaType, meta.target, injection.ComponentInternalScope))
		}
	}

	topics := getPeerTopics(featureHelper.feature)
	if len(topics) > 0 {
		for _, topic := range topics {
			components = append(components, injection.NewComponent(topic, injection.ComponentNoneScope))
		}
	}

	// brokers
	brokerMetaList := getBrokerMetaList(featureHelper.feature)
	if len(brokerMetaList) > 0 {
		for _, meta := range brokerMetaList {
			components = append(components, injection.NewComponentByType(meta.metaType, meta.target, injection.ComponentInternalScope))
		}
	}

	handlers := getSerlvetHandlers(featureHelper.feature)
	if len(handlers) > 0 {
		for _, handler := range handlers {
			components = append(components, injection.NewComponent(handler, injection.ComponentNoneScope))
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

func (module *stdFeatureModule) Defer(ctx context.Context) error {
	featureHelper := module.featureHelper

	// init brokers
	metaList := getBrokerMetaList(featureHelper.feature)
	if len(metaList) > 0 {
		for _, meta := range metaList {
			meta.target.InitBroker(module.featureHelper)
		}
	}
	//

	// attach repositories
	repository := getRepository(featureHelper.feature)
	if repository == nil {
		return nil
	}

	if isTempDB(featureHelper.feature) {
		return module.RepositoryManager.AttachTempModule(module)
	}
	return module.RepositoryManager.AttachBaseModule(module)
	//
}

var _ = (bootstrap.DestroyModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) Destroy() {
	// detach repositories
	featureHelper := module.featureHelper
	repository := getRepository(featureHelper.feature)
	if repository == nil {
		return
	}

	if isTempDB(featureHelper.feature) {
		module.RepositoryManager.DetachTempModule(module)
	} else {
		module.RepositoryManager.DetachBaseModule(module)
	}
	//
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
