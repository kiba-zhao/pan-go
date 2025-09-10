package feature

import "pan/lib/injection"

type SubModuleNewFunc[T any] func(parentModule T) interface{}

func NewSubModule[Module any, ParentModule any](module Module, parentModule ParentModule) *SubModule[Module, ParentModule] {
	subModule := &SubModule[Module, ParentModule]{}
	subModule.specifyModule = module
	subModule.parentModule = parentModule
	return subModule
}

type SubModule[Module any, ParentModule any] struct {
	specifyModule Module
	parentModule  ParentModule
}

var _ = (injection.ComponentStoreProvider)((*SubModule[any, any])(nil))

func (s *SubModule[Module, ParentModule]) ComponentStore() injection.ComponentStore {
	parentModule := s.ParentModule()
	if provider, ok := any(parentModule).(injection.ComponentStoreProvider); ok {
		return provider.ComponentStore()
	}
	return nil
}

func (s *SubModule[Module, ParentModule]) ParentModule() ParentModule {
	return s.parentModule
}
