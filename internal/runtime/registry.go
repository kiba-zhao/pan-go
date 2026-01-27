/*
# Registry

The registry is a thread-safe data structure that allows modules to be registered
and accessed by other modules. It provides a simple and efficient way to manage
and access modules in the engine.
*/
package runtime

import (
	"errors"
	"iter"
	"reflect"
	"slices"
	"sync"
)

// TraverseFunc is a function that takes a module and returns an error
type TraverseFunc[T any] func(module T) error

// Registry is an interface for registering modules
type Registry interface {

	// Count returns the number of modules registered under the specified type.
	// It acquires a read lock to ensure thread-safe access to the registry's modules.
	Count(t reflect.Type) int

	// Append adds a module to the registry under the specified types.
	// If a type is not provided, the type of the module is used.
	// It acquires a write lock to ensure thread-safe access to the registry's modules.
	// Returns ErrModuleType if the module does not implement the specified type(s).
	Append(module interface{}, types ...reflect.Type) error

	// Modules returns a slice of all modules registered in the registry.
	Modules() []interface{}

	// ModulesByType returns a slice of all modules registered in the registry
	// under the specified type, and a boolean indicating whether any modules
	// were found.
	ModulesByType(t reflect.Type) ([]interface{}, bool)

	// Traverse traverses all modules registered in the registry,
	// and calls the specified function on each module.
	// If the function returns an error, traversal stops and the error is returned.
	// It acquires a read lock to ensure thread-safe access to the registry's modules.
	Traverse(f TraverseFunc[interface{}]) error

	// TraverseByType traverses all modules registered in the registry
	// under the specified type, and calls the specified function on each module.
	// If the function returns an error, traversal stops and the error is returned.
	// It acquires a read lock to ensure thread-safe access to the registry's modules.
	TraverseByType(f TraverseFunc[interface{}], t reflect.Type) error
}

var ErrModuleType = errors.New("[engine:Registry] Module Error: wrong type")

type registryImpl struct {
	modules map[reflect.Type][]interface{}
	rw      sync.RWMutex
}

var _ = (Registry)((*registryImpl)(nil))

// NewRegistry creates a new registry.
//
// The registry is a central location for all modules
// in the engine. It provides methods to add, get, and
// traverse modules.
//
// The returned registry is safe for concurrent access.
func NewRegistry() Registry {
	registry := &registryImpl{}
	registry.modules = make(map[reflect.Type][]interface{})
	return registry
}

// Count returns the number of modules registered under the specified type.
// It acquires a read lock to ensure thread-safe access to the registry's modules.
func (r *registryImpl) Count(t reflect.Type) int {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.modules[t])
}

// Append adds a module to the registry under the specified types.
// If a type is not provided, the type of the module is used.
// It acquires a write lock to ensure thread-safe access to the registry's modules.
// Returns ErrModuleType if the module does not implement the specified type(s).
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

// Modules returns a slice of all modules registered in the registry.
// It acquires a read lock to ensure thread-safe access to the registry's modules.
func (r *registryImpl) Modules() []interface{} {
	r.rw.RLock()
	defer r.rw.RUnlock()

	var modules []interface{}
	for _, m := range r.modules {
		modules = append(modules, m...)
	}
	return modules
}

// ModulesByType returns a slice of all modules registered in the registry
// under the specified type, and a boolean indicating whether any modules
// were found.
//
// It acquires a read lock to ensure thread-safe access to the registry's modules.
func (r *registryImpl) ModulesByType(t reflect.Type) ([]interface{}, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	modules, ok := r.modules[t]
	return modules, ok
}

// Traverse traverses all modules registered in the registry,
// and calls the specified function on each module.
// If the function returns an error, traversal stops and the error is returned.
// It acquires a read lock to ensure thread-safe access to the registry's modules.
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

// TraverseByType traverses all modules registered in the registry under the specified type,
// and calls the specified function on each module.
// If the function returns an error, traversal stops and the error is returned.
// It acquires a read lock to ensure thread-safe access to the registry's modules.
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

// TraverseRegistry traverses all modules in the registry that implement the specified type T,
// and applies the provided function to each module. It utilizes the TraverseByType method
// to ensure only modules of the specified type are traversed. If the function returns an error
// for any module, the traversal stops and the error is returned.
func TraverseRegistry[T any](registry Registry, f TraverseFunc[T]) error {
	t := reflect.TypeFor[T]()

	return registry.TraverseByType(func(module interface{}) error {
		return f(module.(T))
	}, t)
}

// ModulesForType returns a slice of all modules registered in the registry
// that implement the specified type T.
//
// It acquires a read lock to ensure thread-safe access to the registry's modules.
//
// If no modules of the specified type are found, ModulesForType returns nil.
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

// SeqForType returns a sequence of all modules registered in the registry
// that implement the specified type T. It utilizes the ModulesByType method
// to ensure only modules of the specified type are iterated over. If no
// modules of the specified type are found, SeqForType returns an empty sequence.
func SeqForType[T any](registry Registry) iter.Seq[T] {

	t := reflect.TypeFor[T]()

	return func(yield func(T) bool) {
		modules, ok := registry.ModulesByType(t)
		if !ok {
			return
		}

		for _, module := range modules {
			m, ok := module.(T)
			if !ok {
				continue
			}
			if !yield(m) {
				return
			}
		}
	}

}
