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
	GetModuleByName(name string) AppModule
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
				return errors.New("module conflict")
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

func (r *registryImpl) GetModuleByName(name string) AppModule {
	r.rw.RLock()
	defer r.rw.RUnlock()
	for _, m := range r.modules {
		if m.Name() == name {
			return m
		}
	}
	return nil
}
