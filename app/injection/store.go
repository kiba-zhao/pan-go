package injection

import "reflect"

type ComponentStore = map[reflect.Type]interface{}

type ComponentStoreProvider interface {
	ComponentStore() ComponentStore
}

func NewComponentStore() ComponentStore {
	return make(ComponentStore)
}
