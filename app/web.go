package app

import (
	"context"
	"errors"
	"io/fs"
	"net/http"

	"pan/runtime"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

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
	app        WebApp
	appLocker  sync.RWMutex
	registry   runtime.Registry
	locker     sync.RWMutex
	addresses  []string
	sigChan    chan bool
	sigOnce    sync.Once
	hasPending bool
}

func (w *webServer) SetSig(sig bool) {

	if w.hasPending {
		return
	}

	w.sigOnce.Do(func() {
		w.sigChan = make(chan bool, 1)
	})

	w.hasPending = true
	w.sigChan <- sig

}

func (w *webServer) Addresses() []string {
	w.locker.RLock()
	defer w.locker.RUnlock()
	return w.addresses
}

func (w *webServer) OnConfigUpdated(settings AppSettings) {
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
	w.SetSig(true)
	w.locker.Unlock()

	return w.ReloadModules()
}

func (w *webServer) EngineTypes() []reflect.Type {
	return []reflect.Type{
		reflect.TypeFor[WebModule](),
		reflect.TypeFor[WebModuleProvider](),
	}
}

func (w *webServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	w.appLocker.RLock()
	app := w.app
	w.appLocker.RUnlock()

	if app != nil {
		app.ServeHTTP(rw, req)
	} else {
		http.Error(rw, ErrUnavailable.Error(), http.StatusInternalServerError)
	}
}

func (w *webServer) ReloadModules() error {
	w.locker.Lock()
	registry := w.registry
	w.locker.Unlock()
	if registry == nil {
		return ErrUnavailable
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
		w.app = app
	}

	return err
}

func (w *webServer) Ready() error {

	var servers []*http.Server
	w.sigOnce.Do(func() {
		w.sigChan = make(chan bool, 1)
	})

	for {

		sig := <-w.sigChan
		w.locker.Lock()
		w.hasPending = false
		w.locker.Unlock()

		if len(servers) > 0 {
			for _, server := range servers {
				server.Shutdown(context.Background())
			}
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
			go func(s *http.Server) {
				for {
					err := s.ListenAndServe()
					if errors.Is(err, http.ErrServerClosed) {
						break
					}
					time.Sleep(6 * time.Second)
				}

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
