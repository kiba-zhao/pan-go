package servlet

import (
	"context"
	"errors"
	"io"
	libApp "pan/lib/app"
	"sync"
)

const (
	CodeOK            = libApp.CodeOK
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeForbidden     = 403
)

var errServletUnavailable = errors.New("servlet.Servlet Error: Unavailable")
var errServletNotFound = errors.New("servlet.Servlet Error:  Not Found")

type ServletApp = *libApp.App
type ServletRouter = libApp.AppHandleGroup
type ServletAction = libApp.RequestName

var ServletDoContext = []byte("ServletDoContext")

type Servlet interface {
	Do(context.Context, ServletAction, io.Reader) (io.ReadCloser, error)
}

type stdServlet struct {
	app   ServletApp
	appRW sync.RWMutex
}

var _ = (Servlet)((*stdServlet)(nil))

func (servlet *stdServlet) Do(ctx context.Context, action ServletAction, reader io.Reader) (io.ReadCloser, error) {
	app := getApp(servlet)
	if app == nil {
		return nil, errServletUnavailable
	}

	req := libApp.NewRequest(action, reader)
	appCtx := &libApp.Context{}
	libApp.InitContextWithRequest(appCtx, req)
	appCtx.Set(ServletDoContext, ctx)

	err := app.Run(appCtx, nil)
	defer releaseReader(reader)
	if err != nil {
		return nil, err
	}

	if appCtx.Code() < 0 {
		return nil, errServletNotFound
	}

	if appCtx.Code() != libApp.CodeOK {
		return nil, appCtx.Err()
	}

	return appCtx, nil
}

func getApp(servlet *stdServlet) ServletApp {
	servlet.appRW.RLock()
	defer servlet.appRW.RUnlock()
	return servlet.app
}

func setApp(servlet *stdServlet, app ServletApp) {
	servlet.appRW.Lock()
	defer servlet.appRW.Unlock()
	servlet.app = app
}

func releaseReader(reader io.Reader) {
	if closer, ok := reader.(io.Closer); ok {
		closer.Close()
	}
}
