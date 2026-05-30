package app

import (
	"errors"
	"io"
	"pan/internal/servlet"
	"sync"
)

var ErrAppletUnavailable = errors.New("net.Applet Error: Applet Unavailable")
var ErrAppletCommandNotFound = errors.New("net.Applet Error: Command Not Found")

type CommandName = servlet.RequestName
type AppServletContext = *servlet.Context
type AppServlet = *servlet.Servlet[AppServletContext]
type AppServletNext = servlet.Next
type AppServletRouter = servlet.HandleGroup[AppServletContext]

type Applet interface {
	Exec(name CommandName, reader io.Reader) (io.ReadCloser, error)
}

type stdApplet struct {
	appServlet   AppServlet
	appServletRW sync.RWMutex
}

var _ = (Applet)((*stdApplet)(nil))

func (applet *stdApplet) Exec(name CommandName, reader io.Reader) (io.ReadCloser, error) {
	applet.appServletRW.RLock()
	appServlet := applet.appServlet
	applet.appServletRW.RUnlock()

	if appServlet == nil {
		return nil, ErrAppletUnavailable
	}

	req := servlet.NewRequest(name, reader)
	appCtx := &servlet.Context{}
	servlet.InitContextWithRequest(appCtx, req)

	err := appServlet.Run(appCtx, nil)
	if err != nil {
		return nil, err
	}

	if appCtx.Code() < 0 {
		return nil, ErrAppletCommandNotFound
	}

	if appCtx.Code() != servlet.CodeOK {
		return nil, appCtx.Err()
	}
	return appCtx, nil
}

func (applet *stdApplet) SetupAppServlet(appServlet AppServlet) {
	applet.appServletRW.Lock()
	defer applet.appServletRW.Unlock()
	applet.appServlet = appServlet
}
