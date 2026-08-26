//go:build !(android || ios)

package settings

import (
	"context"
	"pan/pkg/bootstrap"
	"pan/pkg/injection"
	"pan/pkg/module"
	"pan/pkg/web"
	"path"
)

func init() {
	subModuleNewFuncArray = append(subModuleNewFuncArray, newHostModule)
}

func newHostModule(m *stdModule) interface{} {
	deviceInfoCtrl := &DeviceInfoController{}
	deviceNetworkCtrl := &DeviceNetworkController{}
	webHostCtrl := &WebHostController{}
	webHostSvc := &WebHostService{}

	hostModule := &stdHostSettingsModule{}
	hostModule.SubModule = module.NewSubModule(hostModule, m)
	hostModule.deviceInfoCtrl = deviceInfoCtrl
	hostModule.deviceNetworkCtrl = deviceNetworkCtrl
	hostModule.webHostCtrl = webHostCtrl
	hostModule.webHostSvc = webHostSvc
	webHostSvc.Trigger = hostModule

	return hostModule
}

type stdHostSettingsModule struct {
	WebConfigurer web.WebConfigurer
	Configurer    SettingsConfigurer

	*module.SubModule[*stdHostSettingsModule, *stdModule]

	deviceInfoCtrl    *DeviceInfoController
	deviceNetworkCtrl *DeviceNetworkController
	webHostCtrl       *WebHostController
	webHostSvc        *WebHostService
}

var _ = (web.WebAppModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) SetupToWeb(app web.WebApp) error {
	route := app.Group(path.Join(web.WEB_API_PATH, SettingsModuleName))
	err := m.deviceInfoCtrl.SetupToWeb(route)
	if err == nil {
		err = m.deviceNetworkCtrl.SetupToWeb(route)
	}
	if err == nil {
		err = m.webHostCtrl.SetupToWeb(route)
	}
	return err
}

var _ = (bootstrap.DeferModule)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) Defer(ctx context.Context) error {
	webHost, err := m.loadWebHost()
	if err == nil {
		err = m.configure(webHost)
	}
	return err
}

var _ = (injection.ComponentProvider)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(m, injection.ComponentInternalScope),
		// controller
		injection.NewComponent(m.deviceInfoCtrl, injection.ComponentNoneScope),
		injection.NewComponent(m.deviceNetworkCtrl, injection.ComponentNoneScope),
		injection.NewComponent(m.webHostCtrl, injection.ComponentNoneScope),
		// service
		injection.NewComponent(m.webHostSvc, injection.ComponentInternalScope),
	}
}

var _ = (WebHostChangedTrigger)((*stdHostSettingsModule)(nil))

func (m *stdHostSettingsModule) OnWebHostChanged(webHost WebHost) {
	m.configure(webHost)
}

func (m *stdHostSettingsModule) configure(webHost WebHost) error {
	return m.configureWeb(&webHost)
}

func (m *stdHostSettingsModule) configureWeb(webHost *WebHost) error {
	webConfig := newWebConfig(webHost, m.Configurer.Config())
	return m.WebConfigurer.Configure(webConfig)
}

func (m *stdHostSettingsModule) loadWebHost() (WebHost, error) {
	webHost, err := m.webHostSvc.Load()
	if err != nil {
		m.ParentModule().logger.Error("settings.HostSettingsModule", "loadWebHost Error: "+err.Error())
	}
	return webHost, err
}
