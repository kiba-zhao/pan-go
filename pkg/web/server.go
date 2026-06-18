package web

import (
	"context"
	"net/http"
	"pan/pkg/log"
	"sync"
)

type stdWebServer struct {
	logger log.Logger

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	app    WebApp
	appRW  sync.RWMutex
	addr   string
	addrRW sync.RWMutex
}

func (ws *stdWebServer) Setup(cfg WebConfig) {
	ws.logger.Debug("web.WebServer", "Setup")

	ws.addrRW.Lock()
	defer ws.addrRW.Unlock()

	addr := cfg.Addr()

	if ws.addr == addr {
		return
	}
	ws.addr = addr

	ws.reloadLock.Lock()
	defer ws.reloadLock.Unlock()
	if !ws.reload {
		ws.reload = true
		ws.reloadChan <- struct{}{}
	}

}

func (ws *stdWebServer) SetupApp(app WebApp) {
	ws.logger.Info("web.WebServer", "SetupApp")

	ws.appRW.Lock()
	defer ws.appRW.Unlock()
	if ws.app == app {
		return
	}
	ws.app = app

	ws.reloadLock.Lock()
	defer ws.reloadLock.Unlock()
	if !ws.reload {
		ws.reload = true
		ws.reloadChan <- struct{}{}
	}
}

func (ws *stdWebServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "WebApp is not available", http.StatusInternalServerError)
}

func (ws *stdWebServer) ListenAndServe(ctx context.Context) error {
	ws.logger.Debug("web.WebServer", "ListenAndServe begin")
	defer ws.logger.Debug("web.WebServer", "ListenAndServe end")

	var err error
	var closed bool

	var wg sync.WaitGroup
	var server *http.Server
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-ws.reloadChan:
			ws.reloadLock.Lock()
			ws.reload = false
			ws.reloadLock.Unlock()
		}

		if server != nil {
			server.Shutdown(context.Background())
			wg.Wait()
		}

		if closed {
			break
		}

		ws.addrRW.RLock()
		addr := ws.addr
		ws.addrRW.RUnlock()
		if len(addr) <= 0 {
			continue
		}

		ws.appRW.RLock()
		app := ws.app
		ws.appRW.RUnlock()
		var handler http.Handler
		if app != nil {
			handler = app
		} else {
			handler = ws
		}

		httpServer := &http.Server{
			Addr:    addr,
			Handler: handler,
		}

		server = httpServer
		wg.Add(1)
		go func(s *http.Server) {
			defer wg.Done()
			err = s.ListenAndServe()
			if err != nil {
				ws.logger.Error("web.WebServer", "http.Server.ListenAndServe Error: "+err.Error())
			} else {
				ws.logger.Info("web.WebServer", "http.Server.ListenAndServe Success: "+s.Addr)
			}
		}(httpServer)

	}

	return err
}
