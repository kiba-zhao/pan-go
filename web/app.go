package web

import "github.com/gin-gonic/gin"

type App struct {
	*gin.Engine
}

// Mount ...
func (app *App) Mount(modules ...Module) {
	for _, m := range modules {
		m.SetupToWeb(app)
	}
}

// NewApp ...
func NewApp() *App {
	app := new(App)
	app.Engine = gin.New()

	return app
}
