package web

import (
	"context"
	"net/http"
	"pan/lib/log"
	"slices"
	"sync"
)

type stdWebServer struct {
	logger log.Logger

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	app     WebApp
	appRW   sync.RWMutex
	addrs   []string
	addrsRW sync.RWMutex
}

func (ws *stdWebServer) Addrs() []string {
	ws.addrsRW.RLock()
	defer ws.addrsRW.RUnlock()
	return ws.addrs
}

func (ws *stdWebServer) SetAddrs(addrs []string) {
	ws.logger.Debug("WebServer", "SetAddrs")

	ws.addrsRW.Lock()
	defer ws.addrsRW.Unlock()
	if slices.Equal(ws.addrs, addrs) {
		return
	}
	ws.addrs = addrs
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
	var servers []*http.Server
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

		if len(servers) > 0 {
			for _, server := range servers {
				server.Shutdown(context.Background())
			}
			wg.Wait()
		}

		if closed {
			break
		}

		addrs := ws.Addrs()
		if len(addrs) <= 0 {
			continue
		}

		var handler http.Handler
		if app := ws.WebApp(); app != nil {
			handler = app
		} else {
			handler = ws
		}

		servers = make([]*http.Server, 0)
		for _, address := range addrs {
			httpServer := &http.Server{
				Addr:    address,
				Handler: handler,
			}

			servers = append(servers, httpServer)
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

	}

	return err
}
