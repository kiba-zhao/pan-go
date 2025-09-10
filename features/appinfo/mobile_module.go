//go:build android || ios

package appinfo

import (
	"pan/lib/feature"
	"sync"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newMobileModule)
}

func newMobileModule(module *stdModule) interface{} {
	mobileModule := &stdMobileModule{}
	mobileModule.SubModule = feature.NewSubModule(mobileModule, module)

	return mobileModule
}

type stdMobileModule struct {
	*feature.SubModule[*stdMobileModule, *stdModule]

	serlvets     []feature.ServletHandler
	serlvetsOnce sync.Once
}

var _ = (feature.ServletSubModule)((*stdMobileModule)(nil))

func (m *stdMobileModule) ServletRouteName() []byte {
	return []byte(ModuleName)
}

var _ = (feature.ServletHandlerProvider)((*stdMobileModule)(nil))

func (m *stdMobileModule) ServletHandlers() []feature.ServletHandler {
	m.serlvetsOnce.Do(func() {
		m.serlvets = []feature.ServletHandler{
			&QRCodeServlet{},
		}
	})
	return m.serlvets
}
