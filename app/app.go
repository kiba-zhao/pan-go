package app

type App struct {
	settings *Settings
	config   *Config
	web      *WebApp
}

func New() *App {
	app := new(App)
	app.settings = NewSettings()
	app.config = NewConfig(app)
	app.web = NewWebApp(app)
	return app
}

func (app *App) Init() error {
	err := app.settings.init()
	if err != nil {
		return err
	}
	err = app.config.init()
	if err != nil {
		return err
	}
	err = app.config.Load()
	return err
}

func (app *App) Mount(modules ...interface{}) {
	for _, m := range modules {
		if wm, ok := m.(WebModule); ok {
			app.web.Mount(wm)
		}
	}
}

func (app *App) Run() error {
	return app.web.Run()
}
