//go:build !(android || ios)

package appinfo

import (
	"pan/lib/feature"
	"sync"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(module *stdModule) interface{} {
	hostModule := &stdHostModule{}
	hostModule.SubModule = feature.NewSubModule(hostModule, module)
	return hostModule
}

type stdHostModule struct {
	*feature.SubModule[*stdHostModule, *stdModule]

	controllers     []feature.WebController
	controllersOnce sync.Once
}

var _ = (feature.WebSubModule)((*stdHostModule)(nil))

func (m *stdHostModule) WebRouteName() string {
	return ModuleName
}

var _ = (feature.WebControllerProvider)((*stdHostModule)(nil))

func (m *stdHostModule) WebControllers() []feature.WebController {
	m.controllersOnce.Do(func() {
		m.controllers = []feature.WebController{
			&QRCodeController{},
		}
	})
	return m.controllers
}
