package web

import (
	"context"
	"net/http"
	"pan/lib/log"
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

func (ws *stdWebServer) Addr() string {
	ws.addrRW.RLock()
	defer ws.addrRW.RUnlock()
	return ws.addr
}

func (ws *stdWebServer) SetAddr(addr string) {
	ws.logger.Debug("WebServer", "SetAddr")

	ws.addrRW.Lock()
	defer ws.addrRW.Unlock()
	if ws.addr == addr {
		return
	}
	ws.addr = addr
	ws.Reload()
}

func (ws *stdWebServer) Reload() {
	ws.logger.Debug("WebServer", "Reload")

	ws.reloadLock.Lock()
	defer ws.reloadLock.Unlock()
	if ws.reload {
		return
	}

	ws.reload = true
	ws.reloadChan <- struct{}{}
}

func (ws *stdWebServer) WebApp() WebApp {
	ws.appRW.RLock()
	defer ws.appRW.RUnlock()
	return ws.app
}

func (ws *stdWebServer) SetWebApp(app WebApp) {
	ws.logger.Info("WebServer", "SetWebApp")

	ws.appRW.Lock()
	defer ws.appRW.Unlock()
	if ws.app == app {
		return
	}
	ws.app = app
	ws.Reload()
}

func (ws *stdWebServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "WebApp is not available", http.StatusInternalServerError)
}

func (ws *stdWebServer) ListenAndServe(ctx context.Context) error {
	ws.logger.Debug("WebServer", "ListenAndServe begin")
	defer ws.logger.Debug("WebServer", "ListenAndServe end")

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

		addr := ws.Addr()
		if len(addr) <= 0 {
			continue
		}

		var handler http.Handler
		if app := ws.WebApp(); app != nil {
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
				ws.logger.Error("WebServer", "http.Server.ListenAndServe Error: "+err.Error())
			} else {
				ws.logger.Info("WebServer", "http.Server.ListenAndServe Success: "+s.Addr)
			}
		}(httpServer)

	}

	return err
}
