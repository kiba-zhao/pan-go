package extfs

import (
	"embed"
	"pan/core"
)

//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production --base /extfs/
//go:embed web/dist
var embedFS embed.FS

type Module struct {
	*core.BrowserRouteWebModule
	cfg core.Config
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
	return "extfs"
}

func (m *Module) Desc() string {
	return "ExtFS Module"
}

// TODO: SetupToWeb
func (m *Module) SetupToWeb(router core.WebRouter) {
	// TODO: Dependency Injection
	// Mount Controllers

}

// TODO: OnInitConfig
func (m *Module) OnInitConfig(cfg core.Config) error {
	settings := defaultSettings()
	cfg.Init(settings)
	m.cfg = cfg
	return nil
}

// TODO: HasWeb
func (m *Module) HasWeb() bool {
	return true
}

// TODO: EnabledModule
func (m *Module) Enabled() bool {
	var settings Settings
	err := m.cfg.Load(&settings)
	if err != nil {
		// TODO: log error
		return false
	}
	return settings.Enabled
}

// TODO: SetEnable
func (m *Module) SetEnable(enable bool) error {
	var settings Settings
	err := m.cfg.Load(&settings)
	if err != nil {
		return err
	}
	settings.Enabled = enable
	return m.cfg.Sync(&settings)
}

// TODO: ReadOnlyModule
func (m *Module) ReadOnly() bool {
	return false
}

// initForApp
func (m *Module) InitForApp() error {

	// TODO: trigger tasks
	return nil
}
