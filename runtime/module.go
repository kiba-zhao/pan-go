package runtime

type module struct {
	modules []interface{}
}

func NewModule(modules ...interface{}) interface{} {
	return &module{
		modules: modules,
	}
}

func (s *module) Modules() []interface{} {
	return s.modules
}
