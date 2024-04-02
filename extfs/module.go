package extfs

import (
	"embed"
	"pan/core"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"pan/extfs/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:generate npm --prefix ./web install
//go:generate npm --prefix ./web run build -- -m production --base /extfs/
//go:embed web/dist
var embedFS embed.FS

type Module struct {
	*core.BrowserRouteWebModule
	cfg      core.Config
	settings *services.SettingsService
	db       *gorm.DB
}

func NewModule() *Module {
	m := new(Module)
	m.BrowserRouteWebModule = core.NewBrowserRouteWebModule(embedFS, "web/dist")

	m.settings = &services.SettingsService{}

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

	api := router.Group("/api")
	var targetCtrl controllers.TargetController
	targetCtrl.TargetService = &services.TargetService{}
	targetCtrl.TargetService.TargetRepo = repositories.NewTargetRepository(m.db)
	targetCtrl.Init(api)

	var diskfileCtrl controllers.DiskFileController
	diskfileCtrl.DiskFileService = &services.DiskFileService{}
	diskfileCtrl.Init(api)

}

// TODO: OnInitConfig
func (m *Module) OnInitConfig(cfg core.Config) error {

	// intialize settings
	settings := &Settings{}
	settings.init()
	cfg.Init(settings)
	m.cfg = cfg

	// set settings service properties
	m.settings.Settings = &settings.Settings
	m.settings.Config = cfg

	// create db
	db, err := gorm.Open(sqlite.Open(settings.DBFilePath), &gorm.Config{})
	if err == nil {
		m.db = db
		err = m.InitDB()
	}
	return err
}

// TODO: HasWeb
func (m *Module) HasWeb() bool {
	return true
}

// TODO: EnabledModule
func (m *Module) Enabled() bool {
	settings := &Settings{}
	err := m.cfg.Marshal(settings)
	if err != nil {
		// TODO: log error
		return false
	}
	return settings.Enabled
}

// TODO: SetEnable
func (m *Module) SetEnable(enable bool) error {

	return m.cfg.Sync(func(marshal core.ConfigMarshalHandle) (interface{}, error) {
		settings := &Settings{}
		err := marshal(settings)
		if err != nil {
			return nil, err
		}
		settings.Enabled = enable
		return settings, nil
	})
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

func (m *Module) InitDB() error {
	err := m.db.AutoMigrate(&models.Target{})
	return err
}
