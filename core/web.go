package core

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebRouter = gin.IRouter
type WebContext = *gin.Context

var WrapHandleFunc = gin.WrapF
var WrapHandler = gin.WrapH

type WebModule interface {
	AppModule
	SetupToWeb(app WebRouter)
	HasWeb() bool
}

type WebApp struct {
	*gin.Engine
	settings *Settings
}

func NewWebApp(settings *Settings) *WebApp {
	web := new(WebApp)
	web.Engine = gin.New()
	web.settings = settings

	return web
}

// Mount ...
func (web *WebApp) Mount(modules ...WebModule) {
	for _, m := range modules {
		name := m.Name()
		rg := web.Group(name)
		m.SetupToWeb(rg)
	}
}

func (web *WebApp) Run() error {
	settings := web.settings

	if strings.Trim(settings.AppModule, " ") != "" {
		location := "/" + settings.AppModule
		rootRedirect := func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, location)
		}
		web.GET("/", rootRedirect)
		web.HEAD("/", rootRedirect)
	}

	return web.Engine.Run(settings.WebHost + ":" + strconv.Itoa(settings.WebPort))
}
