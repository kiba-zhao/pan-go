package app

import (
	"embed"
	"pan/app/controllers"
	"pan/app/services"
	"pan/core"
)

//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production --base /app/
//go:embed web/dist
var embedFS embed.FS

type Module struct {
	*core.BrowserRouteWebModule
	registry core.Registry
}

func NewModule() *Module {

	m := new(Module)

	m.BrowserRouteWebModule = core.NewBrowserRouteWebModule(embedFS, "web/dist")

	return m
}

func (m *Module) Avatar() string {
	return ""
}

func (m *Module) Name() string {
	return "app"
}

func (m *Module) Desc() string {
	return "App Module"
}

func (m *Module) SetupToWeb(router core.WebRouter) {

	// TODO: Dependency Injection
	// Mount Controllers
	ctrlRouter := router.Group("/api")
	var ctrl controllers.ModuleController
	ctrl.ModuleService = &services.ModuleService{}
	ctrl.ModuleService.Registry = m.registry
	ctrl.Init(ctrlRouter)

}

func (m *Module) HasWeb() bool {
	return true
}

func (m *Module) OnAddRegistry(registry core.Registry) error {
	m.registry = registry
	return nil
}
