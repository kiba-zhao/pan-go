package feature

import (
	"pan/internal/injection"
	"pan/internal/runtime"
	"sync"
)

func NewWithSubModules(module interface{}, subModules ...interface{}) *FeatureModule {
	featureModule := &FeatureModule{}
	featureModule.subModules = subModules
	return featureModule
}

func New[T any](module T, newFuncArr ...SubModuleNewFunc[T]) *FeatureModule {
	featureModule := &FeatureModule{}
	featureModule.subModules = make([]interface{}, 0)
	if len(newFuncArr) > 0 {
		for _, newFunc := range newFuncArr {
			featureModule.subModules = append(featureModule.subModules, newFunc(module))
		}
	}
	return featureModule
}

type FeatureModule struct {
	componentStore     injection.ComponentStore
	componentStoreOnce sync.Once

	subModules []interface{}
}

var _ = (injection.ComponentStoreProvider)((*FeatureModule)(nil))

func (module *FeatureModule) ComponentStore() injection.ComponentStore {
	module.componentStoreOnce.Do(func() {
		module.componentStore = injection.NewComponentStore()
	})
	return module.componentStore
}

var _ = (runtime.ProviderModule)((*FeatureModule)(nil))

func (module *FeatureModule) Modules() []interface{} {
	return module.subModules
}
