// Define injection store
package injection

import "reflect"

type ComponentStore = map[reflect.Type]interface{}

type ComponentStoreProvider interface {
	// ComponentStore returns the ComponentStore.
	//
	// A ComponentStore is a map that uses a reflect.Type as the key and stores
	// the value of the component that matches the type.
	//
	ComponentStore() ComponentStore
}

// NewComponentStore creates a new ComponentStore.
//
// A ComponentStore is a map that uses a reflect.Type as the key and stores
// the value of the component that matches the type.
//
// The ComponentStore is used to store components that are injected into
// other components.
func NewComponentStore() ComponentStore {
	return make(ComponentStore)
}
