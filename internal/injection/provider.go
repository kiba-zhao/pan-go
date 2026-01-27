package injection

type ComponentProvider interface {
	// Components returns a slice of Component instances
	// that this provider is responsible for managing.
	Components() []Component
}

// NewComponentProvider returns a ComponentProvider that manages the given
// components.
func NewComponentProvider(components ...Component) ComponentProvider {
	provider := &stdComponentProvider{}
	provider.components = components
	return provider
}

type stdComponentProvider struct {
	components []Component
}

func (provider *stdComponentProvider) Components() []Component {
	return provider.components
}

// NewStoreComponentProvider returns a ComponentProvider that manages the given
// components and uses the given store to satisfy dependencies.
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
