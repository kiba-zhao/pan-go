//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production
package app

import (
	"embed"
	"io/fs"
	"net/http"
	"pan/app/controllers"
	"pan/app/services"
	"pan/core"
)

//go:embed web/dist
var embedFS embed.FS

type Module struct {
	settings *Settings
	cfg      core.Config
	registry core.Registry
}

func NewModule() *Module {

	m := new(Module)
	m.settings = &Settings{}
	initSettings(m.settings)

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

func (m *Module) Settings() interface{} {
	return m.settings
}

func (m *Module) OnInitConfig(cfg core.Config) error {
	m.cfg = cfg
	return nil
}

func (m *Module) SetupToWeb(router core.WebRouter) {

	// Mount Controllers
	var ctrl controllers.ModuleController
	ctrl.ModuleService = &services.ModuleService{}
	ctrl.ModuleService.Registry = m.registry
	ctrl.Init(router)

	// Mount Static Files
	dist, err := fs.Sub(embedFS, "web/dist")
	if err != nil {
		panic(err)
	}

	router.StaticFileFS("/", "./", http.FS(dist))
	fs.WalkDir(dist, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && d.IsDir() == false {
			router.StaticFileFS(path, path, http.FS(dist))
		}
		return err
	})

}

func (m *Module) HasWeb() bool {
	return true
}

func (m *Module) OnAddRegistry(registry core.Registry) error {
	m.registry = registry
	return nil
}
