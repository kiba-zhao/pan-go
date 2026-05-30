package appinfo

import (
	"pan/internal/injection"
	"pan/internal/module"
)

var subModuleNewFuncArray []module.SubModuleNewFunc[*stdModule]

func New() interface{} {

	m := &stdModule{}
	m.BaseModule = module.New(m, subModuleNewFuncArray...)

	return m
}

const (
	ModuleName = "app-info"
)

type stdModule struct {
	*module.BaseModule
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdModule) Components() []injection.Component {
	return []injection.Component{
		// service
		injection.NewComponent(&QRCodeService{}, injection.ComponentInternalScope),
	}
}
