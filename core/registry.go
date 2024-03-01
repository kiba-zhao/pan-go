package core

import (
	"errors"
	"slices"
	"sync"
)

type RegistryModule interface {
	OnAddRegistry(registry Registry) error
}

type Registry interface {
	AddModule(module AppModule) error
	GetModules() []AppModule
}

type registryImpl struct {
	modules []AppModule
	rw      sync.RWMutex
}

func (r *registryImpl) AddModule(module AppModule) error {

	registryModule, ok := module.(RegistryModule)
	if ok {
		err := registryModule.OnAddRegistry(r)
		if err != nil {
			return err
		}
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	if len(r.modules) > 0 {
		for _, m := range r.modules {
			if m.Name() == module.Name() {
				return errors.New("Module Conflict")
			}
		}
	}
	r.modules = append(r.modules, module)

	return nil
}

func (r *registryImpl) GetModules() []AppModule {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return slices.Clone(r.modules)
}
