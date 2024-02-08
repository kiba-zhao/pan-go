package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type WebRouter = gin.IRouter
type WebContext = *gin.Context

type WebModule interface {
	WebName() string
	SetupToWeb(app WebRouter)
}

type WebApp struct {
	*gin.Engine
	app *App
}

// Mount ...
func (web *WebApp) Mount(modules ...WebModule) {
	for _, m := range modules {
		name := m.WebName()
		if len(name) > 0 {
			rg := web.Group(name)
			m.SetupToWeb(rg)
		} else {
			m.SetupToWeb(web)
		}
	}
}

func (web *WebApp) Run() error {
	settings := web.app.settings
	return web.Engine.Run(settings.WebHost + ":" + strconv.Itoa(settings.WebPort))
}

func NewWebApp(app *App) *WebApp {
	web := new(WebApp)
	web.Engine = gin.New()
	web.app = app
	return web
}
