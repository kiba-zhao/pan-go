package core

import "path"

type AppModule interface {
	Avatar() string
	Name() string
	Desc() string
}

type App struct {
	settings *Settings
	config   *configImpl
	web      *WebApp
	registry Registry
}

func New() *App {
	app := new(App)
	app.settings = &Settings{}
	app.config = &configImpl{}
	app.web = NewWebApp(app.settings)
	app.registry = &registryImpl{}
	return app
}

func (app *App) Init() error {

	err := app.settings.init()
	if err != nil {
		return err
	}

	configPath := path.Join(app.settings.AppRoot(), app.settings.AppName()+".toml")
	err = app.config.init(configPath, app.settings)

	return err
}

func (app *App) Mount(modules ...interface{}) {
	for _, module := range modules {

		m, ok := module.(AppModule)
		if ok == false {
			// TODO: Log error
			continue
		}

		cm, ok := module.(ConfigModule)
		if ok {
			cfg := &configImpl{}
			cfgPath := path.Join(app.settings.AppRoot(), "conf.d", m.Name()+".toml")
			err := cfg.init(cfgPath, cm.Settings())
			if err == nil {
				cm.OnInitConfig(cfg)
			}
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

func (app *App) Run() error {
	return app.web.Run()
}
