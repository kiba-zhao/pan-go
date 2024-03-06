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
	settings *Settings
	cfg      core.Config
}

func NewModule() *Module {
	m := new(Module)
	m.BrowserRouteWebModule = core.NewBrowserRouteWebModule(embedFS, "web/dist")
	m.settings = defaultSettings()

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

// TODO: ConfigModule with settings
func (m *Module) Settings() interface{} {
	return m.settings
}

// TODO: OnInitConfig
func (m *Module) OnInitConfig(cfg core.Config) error {

	m.cfg = cfg
	return nil
}

// TODO: HasWeb
func (m *Module) HasWeb() bool {
	return true
}

// TODO: EnabledModule
func (m *Module) Enabled() bool {
	return m.settings.Enabled
}

// TODO: SetEnable
func (m *Module) SetEnable(enable bool) error {
	m.settings.Enabled = enable
	return m.cfg.Save()
}

// TODO: ReadOnlyModule
func (m *Module) ReadOnly() bool {
	return false
}
