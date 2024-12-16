package injection

type ComponentProvider interface {
	Components() []Component
}

type simpleComponentProvider struct {
	components []Component
}

func NewComponentProvider(components ...Component) ComponentProvider {
	provider := &simpleComponentProvider{}
	provider.components = components
	return provider
}

func (s *simpleComponentProvider) Components() []Component {
	return s.components
}
