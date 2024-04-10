package runtime

type simpleModule struct {
	modules []interface{}
}

func NewSimpleModule(modules ...interface{}) interface{} {
	return &simpleModule{
		modules: modules,
	}
}

func (s *simpleModule) GetModules() []interface{} {
	return s.modules
}
