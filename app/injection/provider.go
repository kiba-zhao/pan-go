package injection

type ComponentProvider interface {
	Components() []Component
}

func NewComponentProvider(components ...Component) ComponentProvider {
	provider := &simpleComponentProvider{}
	provider.components = components
	return provider
}

type simpleComponentProvider struct {
	components []Component
}

func (provider *simpleComponentProvider) Components() []Component {
	return provider.components
}

func NewStoreComponentProvider(store ComponentStore, components ...Component) ComponentProvider {
	provider := &storeComponentProvider{}
	provider.store = store
	provider.ComponentProvider = NewComponentProvider(components...)
	return provider
}

type storeComponentProvider struct {
	ComponentProvider
	store ComponentStore
}

func (provider *storeComponentProvider) ComponentStore() ComponentStore {
	return provider.store
}
