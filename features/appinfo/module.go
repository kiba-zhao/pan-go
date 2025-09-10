package appinfo

import (
	"pan/lib/feature"
	"pan/lib/injection"
)

var subModuleNewFuncArray []feature.SubModuleNewFunc[*stdModule]

func New() interface{} {

	module := &stdModule{}
	module.FeatureModule = feature.New(module, subModuleNewFuncArray...)

	return module
}

const (
	ModuleName = "app-info"
)

type stdModule struct {
	*feature.FeatureModule
}

var _ = (injection.ComponentProvider)((*stdModule)(nil))

func (m *stdModule) Components() []injection.Component {
	return []injection.Component{
		// service
		injection.NewComponent(&QRCodeService{}, injection.ComponentInternalScope),
	}
}
