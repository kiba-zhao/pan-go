package net

import (
	"context"
	"io/fs"
	"net/http"
	"pan/app/constant"

	"pan/app/config"
	"pan/runtime"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type WebApp = *gin.Engine
type WebRouter = gin.IRouter
type WebContext = *gin.Context

func NewWebApp() WebApp {
	return gin.New()
}

type WebModule interface {
	SetupToWeb(WebApp) error
}

type WebModuleProvider interface {
	WebModules() []WebModule
}

type WebController interface {
	SetupToWeb(WebRouter) error
}

type WebScopeModule interface {
	WebScope() string
}

type WebControllerProvider interface {
	WebControllers() []WebController
}

type webServer struct {
	app       WebApp
	appLocker sync.RWMutex
	registry  runtime.Registry
	locker    sync.RWMutex
	addresses []string
	sigChan   chan bool
	sigOnce   sync.Once
	hasSig    bool
}

func (w *webServer) SigChan() chan bool {

	w.sigOnce.Do(func() {
		w.sigChan = make(chan bool, 1)
	})

	return w.sigChan
}

func (w *webServer) SetSig(sig bool) {

	if w.hasSig {
		return
	}

	w.hasSig = true
	w.SigChan() <- sig

}

func (w *webServer) Addresses() []string {
	w.locker.RLock()
	defer w.locker.RUnlock()
	return w.addresses
}

func (w *webServer) OnConfigUpdated(settings config.AppSettings) {
	w.locker.Lock()
	defer w.locker.Unlock()

	if slices.Equal(w.addresses, settings.WebAddress) {
		return
	}

	w.addresses = settings.WebAddress
	w.SetSig(true)
}

func (w *webServer) Init(registry runtime.Registry) error {
	w.locker.Lock()
	w.registry = registry
	w.locker.Unlock()

	return w.ReloadModules()
}

func (w *webServer) Defer() error {
	return w.ReloadModules()
}

func (w *webServer) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebModule](),
		reflect.TypeFor[WebModuleProvider](),
		reflect.TypeFor[WebControllerProvider](),
	}
}

func (w *webServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	w.appLocker.RLock()
	app := w.app
	w.appLocker.RUnlock()

	if app != nil {
		app.ServeHTTP(rw, req)
	} else {
		http.Error(rw, constant.ErrUnavailable.Error(), http.StatusInternalServerError)
	}
}

func (w *webServer) ReloadModules() error {
	w.locker.Lock()
	registry := w.registry
	w.locker.Unlock()
	if registry == nil {
		return constant.ErrUnavailable
	}

	w.appLocker.Lock()
	defer w.appLocker.Unlock()

	app := NewWebApp()

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
	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebControllerProvider) error {
			var router WebRouter
			if scopeModule, ok := module.(WebScopeModule); ok {
				scope := scopeModule.WebScope()
				router = app.Group(scope)
			} else {
				router = app
			}
			for _, wc := range module.WebControllers() {
				perr := wc.SetupToWeb(router)
				if perr != nil {
					return perr
				}
			}
			return nil
		})
	}

	if err == nil {
		w.app = app
	}

	return err
}

func (w *webServer) Ready() error {

	var wg sync.WaitGroup
	var servers []*http.Server
	for {

		sig := <-w.SigChan()
		w.locker.Lock()
		w.hasSig = false
		w.locker.Unlock()

		if len(servers) > 0 {
			for _, server := range servers {
				server.Shutdown(context.Background())
			}
			wg.Wait()
		}

		if !sig {
			break
		}

		addresses := w.Addresses()
		for _, address := range addresses {
			httpServer := &http.Server{
				Addr:    address,
				Handler: w,
			}

			servers = append(servers, httpServer)
			wg.Add(1)
			go func(s *http.Server) {
				defer wg.Done()
				_ = s.ListenAndServe()
				// TODO: echo error into log
			}(httpServer)
		}

	}

	return nil
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
