package web

import (
	"context"
	"errors"
	"net/http"
	"pan/app/config"
	"pan/logger"
	"pan/runtime"
	"reflect"
	"slices"
	"sync"

	"github.com/gin-gonic/gin"
)

var ErrWebModuleUnavailable = errors.New("web.WebModule Error: Unavailable")
var ErrWebModuleInvalidApp = errors.New("web.WebModule Error: Invalid App")

type WebApp = *gin.Engine
type WebRouter = gin.IRouter
type WebContext = *gin.Context

func NewWebApp() WebApp {
	return gin.New()
}

type WebAppModule interface {
	SetupToWeb(WebApp) error
}

type WebAppModuleProvider interface {
	WebAppModules() []WebAppModule
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

func New() interface{} {
	return &webModule{}
}

type webModule struct {
	app        WebApp
	appLocker  sync.RWMutex
	registry   runtime.Registry
	locker     sync.RWMutex
	addresses  []string
	reloadChan chan struct{}
	reloadOnce sync.Once
	needReload bool
}

func (w *webModule) ReloadChan() chan struct{} {

	w.reloadOnce.Do(func() {
		w.reloadChan = make(chan struct{}, 1)
	})

	return w.reloadChan
}

func (w *webModule) OnConfigUpdated(settings config.AppSettings) {
	w.locker.Lock()
	defer w.locker.Unlock()

	if slices.Equal(w.addresses, settings.WebAddress) {
		return
	}

	w.addresses = settings.WebAddress

	// trigger to reload
	if w.needReload {
		return
	}
	w.needReload = true
	w.ReloadChan() <- struct{}{}
}

func (w *webModule) Init(registry runtime.Registry) error {
	w.locker.Lock()
	w.registry = registry
	w.locker.Unlock()

	return w.ReloadModules()
}

func (w *webModule) Defer() error {
	return w.ReloadModules()
}

func (w *webModule) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebAppModule](),
		reflect.TypeFor[WebAppModuleProvider](),
		reflect.TypeFor[WebControllerProvider](),
	}
}

func (w *webModule) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	w.appLocker.RLock()
	app := w.app
	w.appLocker.RUnlock()

	if app != nil {
		app.ServeHTTP(rw, req)
	} else {
		http.Error(rw, ErrWebModuleInvalidApp.Error(), http.StatusInternalServerError)
	}
}

func (w *webModule) ReloadModules() error {
	w.locker.Lock()
	registry := w.registry
	w.locker.Unlock()
	if registry == nil {
		return ErrWebModuleUnavailable
	}

	w.appLocker.Lock()
	defer w.appLocker.Unlock()

	app := NewWebApp()

	err := runtime.TraverseRegistry(registry, func(module WebAppModule) error {
		return module.SetupToWeb(app)
	})
	if err == nil {
		err = runtime.TraverseRegistry(registry, func(module WebAppModuleProvider) error {
			for _, wc := range module.WebAppModules() {
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

func (w *webModule) Ready(ctx context.Context) error {

	var wg sync.WaitGroup
	var servers []*http.Server
	var err error
	closed := false
	for {

		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-w.ReloadChan():
		}
		w.locker.Lock()
		w.needReload = false
		addresses := w.addresses
		w.locker.Unlock()

		if len(servers) > 0 {
			for _, server := range servers {
				server.Shutdown(context.Background())
			}
			wg.Wait()
		}

		if closed {
			break
		}

		for _, address := range addresses {
			httpServer := &http.Server{
				Addr:    address,
				Handler: w,
			}

			servers = append(servers, httpServer)
			wg.Add(1)
			go func(s *http.Server) {
				defer wg.Done()
				err = s.ListenAndServe()
				if err != nil {
					logger.Default().Log(context.Background(), logger.LevelError, "app.web.webModule.Ready Error: %s", err.Error())
				}
				// TODO: echo error into log
			}(httpServer)
		}

	}

	return err
}
