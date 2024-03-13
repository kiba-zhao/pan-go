package core

type AppModule interface {
	Avatar() string
	Name() string
	Desc() string
}

type InitAppModule interface {
	InitForApp() error
}

type App struct {
	settings *Settings
	config   Config
	web      *WebApp
	registry Registry
}

func New() *App {
	app := new(App)
	app.settings = &Settings{}
	app.web = NewWebApp(app.settings)
	app.registry = &registryImpl{}

	cfgPath := generateAppConfigPath(app.settings)
	app.config = newConfig(cfgPath)
	return app
}

func (app *App) Init() error {

	err := app.settings.init()
	if err == nil {
		app.config.Init(app.settings)
	}

	return err
}

func (app *App) Mount(modules ...interface{}) {
	for _, module := range modules {

		m, ok := module.(AppModule)
		if !ok {
			// TODO: Log error
			continue
		}

		cm, ok := module.(ConfigModule)
		if ok {
			cfgPath := generateConfigPath(app.settings, m.Name())
			cfg := newConfig(cfgPath)
			err := cm.OnInitConfig(cfg)
			if err != nil {
				// TODO: Log error
				continue
			}
		}

		ccm, ok := module.(CoreConfigModule)
		if ok {
			err := ccm.OnInitCoreConfig(app.config)
			if err != nil {
				// TODO: Log error
				continue
			}
		}

		err := app.registry.AddModule(m)
		if err != nil {
			// TODO: Log error
			continue
		}

		if wm, ok := m.(WebModule); ok {
			app.web.Mount(wm)
		}

	}
}

func (app *App) Run() (err error) {

	modules := app.registry.GetModules()
	for _, m := range modules {
		if im, ok := m.(InitAppModule); ok {
			err = im.InitForApp()
			if err != nil {
				break
			}
		}
	}

	if err == nil {
		err = app.web.Run()
	}
	return err
}
