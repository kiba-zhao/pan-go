package runtime

import (
	"errors"
	"reflect"
	"slices"
	"sync"
)

type TraverseFunc[T any] func(module T) error

type Registry interface {
	Count(t reflect.Type) int
	Append(module interface{}, types ...reflect.Type) error
	Modules() []interface{}
	ModulesByType(t reflect.Type) ([]interface{}, bool)
	Traverse(f TraverseFunc[interface{}]) error
	TraverseByType(f TraverseFunc[interface{}], t reflect.Type) error
}

var ErrModuleType = errors.New("[engine:Registry] Module Error: wrong type")

type registryImpl struct {
	modules map[reflect.Type][]interface{}
	rw      sync.RWMutex
}

func NewRegistry() Registry {
	registry := &registryImpl{}
	registry.modules = make(map[reflect.Type][]interface{})
	return registry
}

func (r *registryImpl) Count(t reflect.Type) int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.modules[t])
}

func (r *registryImpl) Append(module interface{}, types ...reflect.Type) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	var err error
	t := reflect.TypeOf(module)
	for _, type_ := range types {
		if !(type_.Kind() == reflect.Interface && t.Implements(type_) || type_ == t) {
			err = ErrModuleType
			break
		}
		if modules, ok := r.modules[type_]; ok && !slices.Contains(modules, module) {
			r.modules[type_] = append(modules, module)
			continue
		}
		r.modules[type_] = []interface{}{module}
	}
	return err
}

func (r *registryImpl) Modules() []interface{} {
	r.rw.RLock()
	defer r.rw.RUnlock()

	var modules []interface{}
	for _, m := range r.modules {
		modules = append(modules, m...)
	}
	return modules
}

func (r *registryImpl) ModulesByType(t reflect.Type) ([]interface{}, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	modules, ok := r.modules[t]
	return modules, ok
}

func (r *registryImpl) Traverse(f TraverseFunc[interface{}]) error {

	r.rw.RLock()
	keys := make([]reflect.Type, 0, len(r.modules))
	for k := range r.modules {
		keys = append(keys, k)
	}
	r.rw.RUnlock()

	var err error
	for _, key := range keys {

		r.rw.RLock()
		ms, ok := r.modules[key]
		if !ok || len(ms) == 0 {
			r.rw.RUnlock()
			continue
		}
		modules := slices.Clone(ms)
		r.rw.RUnlock()

		for _, module := range modules {
			err = f(module)
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func (r *registryImpl) TraverseByType(f TraverseFunc[interface{}], t reflect.Type) error {

	r.rw.RLock()
	ms, ok := r.modules[t]

	if !ok || len(ms) == 0 {
		r.rw.RUnlock()
		return nil
	}
	modules := slices.Clone(ms)
	r.rw.RUnlock()

	var err error
	for _, module := range modules {
		err = f(module)
		if err != nil {
			break
		}
	}

	return err
}

func TraverseRegistry[T any](registry Registry, f TraverseFunc[T]) error {
	t := reflect.TypeFor[T]()

	return registry.TraverseByType(func(module interface{}) error {
		return f(module.(T))
	}, t)
}

func ModulesForType[T any](registry Registry) []T {

	t := reflect.TypeFor[T]()
	modules, ok := registry.ModulesByType(t)
	if !ok {
		return nil
	}
	var ts []T
	for _, module := range modules {
		ts = append(ts, module.(T))
	}
	return ts
}
