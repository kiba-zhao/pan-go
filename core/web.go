package core

import (
	"io/fs"
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
	SetupToWeb(router WebRouter)
	HasWeb() bool
}

type AssetsWebModule interface {
	SetupAssetsToWeb(router WebRouter)
}

type noRouteHandlerFunc = gin.HandlerFunc
type NoRouteWebModule interface {
	NoRoute(ctx WebContext)
}

const (
	MIMEHTML = "text/html"
)

func wrapNoRouteHandleFunc(name string, m NoRouteWebModule) noRouteHandlerFunc {
	index := len(name) + 1
	return func(ctx WebContext) {
		path := ctx.Request.URL.Path
		accepted := ctx.NegotiateFormat(MIMEHTML) == MIMEHTML
		if accepted && strings.Index(path, name) == 1 && path[index] == '/' {
			m.NoRoute(ctx)
			return
		}
		ctx.Next()
	}
}

type BrowserRouteWebModule struct {
	distFS   fs.FS
	staticFS http.FileSystem
}

func NewBrowserRouteWebModule(baseFS fs.FS, root string) *BrowserRouteWebModule {
	var distFS fs.FS
	if root == "" {
		distFS = baseFS
	} else {
		subFS, err := fs.Sub(baseFS, root)
		if err != nil {
			panic(err)
		}
		distFS = subFS
	}
	return &BrowserRouteWebModule{
		distFS:   distFS,
		staticFS: http.FS(distFS),
	}
}

func (b *BrowserRouteWebModule) SetupAssetsToWeb(router WebRouter) {

	entries, err := fs.ReadDir(b.distFS, ".")
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() == false {
			router.StaticFileFS(name, name, http.FS(b.distFS))
			continue
		}
		dirFS, err := fs.Sub(b.distFS, name)
		if err != nil {
			panic(err)
		}
		router.StaticFS(name, http.FS(dirFS))
	}
}

// TODO: NoRoute
func (b *BrowserRouteWebModule) NoRoute(ctx WebContext) {
	ctx.FileFromFS("/", b.staticFS)
}

type WebApp struct {
	*gin.Engine
	settings *Settings
	noRoutes []noRouteHandlerFunc
}

func NewWebApp(settings *Settings) *WebApp {
	web := new(WebApp)
	web.Engine = gin.New()
	web.settings = settings
	web.noRoutes = make([]noRouteHandlerFunc, 0)

	return web
}

// Mount ...
func (web *WebApp) Mount(modules ...WebModule) {
	for _, m := range modules {
		name := m.Name()
		rg := web.Group(name)
		m.SetupToWeb(rg)

		am, ok := m.(AssetsWebModule)
		if ok {
			am.SetupAssetsToWeb(rg)
		}

		nrm, ok := m.(NoRouteWebModule)
		if !ok {
			continue
		}
		web.noRoutes = append(web.noRoutes, wrapNoRouteHandleFunc(name, nrm))
	}
}

func (web *WebApp) Run() error {
	settings := web.settings

	if strings.Trim(settings.AppModule, " ") != "" {
		location := "/" + settings.AppModule + "/"
		rootRedirect := func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, location)
		}
		web.GET("/", rootRedirect)
		web.HEAD("/", rootRedirect)
	}

	web.NoRoute(web.noRoutes...)
	return web.Engine.Run(settings.WebHost + ":" + strconv.Itoa(settings.WebPort))
}
