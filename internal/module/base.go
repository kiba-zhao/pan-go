package module

import (
	"pan/internal/injection"
	"pan/internal/runtime"
	"sync"
)

func NewWithSubModules(module interface{}, subModules ...interface{}) *BaseModule {
	featureModule := &BaseModule{}
	featureModule.subModules = subModules
	return featureModule
}

func New[T any](module T, newFuncArr ...SubModuleNewFunc[T]) *BaseModule {
	featureModule := &BaseModule{}
	featureModule.subModules = make([]interface{}, 0)
	if len(newFuncArr) > 0 {
		for _, newFunc := range newFuncArr {
			featureModule.subModules = append(featureModule.subModules, newFunc(module))
		}
	}
	return featureModule
}

type BaseModule struct {
	componentStore     injection.ComponentStore
	componentStoreOnce sync.Once

	subModules []interface{}
}

var _ = (injection.ComponentStoreProvider)((*BaseModule)(nil))

func (module *BaseModule) ComponentStore() injection.ComponentStore {
	module.componentStoreOnce.Do(func() {
		module.componentStore = injection.NewComponentStore()
	})
	return module.componentStore
}

var _ = (runtime.ProviderModule)((*BaseModule)(nil))

func (module *BaseModule) Modules() []interface{} {
	return module.subModules
}
