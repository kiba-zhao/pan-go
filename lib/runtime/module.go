/*
# module

This is a simple Engine Module structure.
It is used to inject the newly created submodule into the registration during Engine Mount.
*/
package runtime

type module struct {
	modules []interface{}
}

// NewModule creates a new instance of a simple engine module.
//
// It takes a variable number of arguments representing the sub-modules to be
// registered in the engine. The sub-modules are returned as a slice of
// interfaces by the Modules method.
func NewModule(modules ...interface{}) interface{} {
	return &module{
		modules: modules,
	}
}

// Modules returns the slice of sub-modules.
//
// This method is part of the Engine's ProviderModule interface, and is used
// to register the sub-modules in the engine during Mount.
func (s *module) Modules() []interface{} {
	return s.modules
}
