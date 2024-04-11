package app

import (
	"io/fs"
	"net/http"
	"pan/runtime"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebApp = *gin.Engine
type WebRouter = gin.IRouter
type WebContext = *gin.Context

func NewWebApp() WebApp {
	return gin.New()
}

type WebModule interface {
	SetupToWeb(app WebApp) error
}

type WebModuleProvider interface {
	WebModules() []WebModule
}

type webServer struct {
	webApp WebApp
	Config AppConfig
}

func (w *webServer) Init(registry runtime.Registry) error {

	// init webApp
	app := NewWebApp()
	w.webApp = app

	err := runtime.TraverseRegistry(registry, func(module WebModule) error {
		return module.SetupToWeb(app)
	})
	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebModuleProvider) error {
			for _, wc := range module.WebModules() {
				perr := wc.SetupToWeb(app)
				if perr != nil {
					return perr
				}
			}
			return nil
		})
	}
	return err
}

func (w *webServer) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebModule](),
		reflect.TypeFor[WebModuleProvider](),
	}
}

func (w *webServer) Components() []Component {
	return []Component{
		NewComponent(w, ComponentNoneScope),
	}
}

func (w *webServer) Ready() error {
	settings, err := w.Config.Read()
	if err == nil {
		err = w.webApp.Run(settings.WebHost + ":" + strconv.Itoa(settings.WebPort))
	}
	return err
}

type webAssets struct {
	basename string
	assets   fs.FS
}

func NewWebAssets(basename string, assets fs.FS) interface{} {
	return &webAssets{
		basename: basename,
		assets:   assets,
	}
}

func (wa *webAssets) SetupToWeb(app WebApp) error {

	hasIndexFile := false
	entries, err := fs.ReadDir(wa.assets, ".")
	if err != nil {
		return err
	}
	var baseRouter WebRouter
	if wa.basename == "" || wa.basename == "/" {
		baseRouter = app
	} else {
		baseRouter = app.Group(wa.basename)
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() {
			baseRouter.StaticFileFS(name, name, http.FS(wa.assets))
			if !hasIndexFile && name == "index.html" {
				hasIndexFile = true
			}
			continue
		}
		dirFS, err := fs.Sub(wa.assets, name)
		if err != nil {
			return err
		}
		baseRouter.StaticFS(name, http.FS(dirFS))
	}

	if hasIndexFile {
		app.NoRoute(wa.noRoute)
	}

	return nil
}

const (
	MIMEHTML = "text/html"
)

func (wa *webAssets) noRoute(ctx WebContext) {
	accepted := ctx.NegotiateFormat(MIMEHTML) == MIMEHTML
	if accepted && strings.HasPrefix(ctx.Request.URL.Path, wa.basename) {
		ctx.FileFromFS("/", http.FS(wa.assets))
		return
	}
	ctx.Next()
}

const (
	CountHeaderName = "X-Total-Count"
)

func SetCountHeaderForWeb(ctx WebContext, total int64) {
	ctx.Header(CountHeaderName, strconv.FormatInt(total, 10))
	ctx.Header("Access-Control-Expose-Headers'", CountHeaderName)
}
