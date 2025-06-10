package serlvet

import (
	"context"
	"errors"
	"io"
	libApp "pan/lib/app"
	"sync"
)

var ErrSerlvetUnavailable = errors.New("serlvet.Serlvet Error: Unavailable")

type SerlvetApp = *libApp.App
type SerlvetRouter = libApp.AppHandleGroup
type SerlvetAction = libApp.RequestName

var SerlvetDoContext = []byte("SerlvetDoContext")

type Serlvet interface {
	Do(context.Context, SerlvetAction, io.Reader) (io.ReadCloser, error)
}

type stdSerlvet struct {
	app   SerlvetApp
	appRW sync.RWMutex
}

var _ = (Serlvet)((*stdSerlvet)(nil))

func (serlvet *stdSerlvet) Do(ctx context.Context, action SerlvetAction, reader io.Reader) (io.ReadCloser, error) {
	app := getApp(serlvet)
	if app == nil {
		return nil, ErrSerlvetUnavailable
	}

	req := libApp.NewRequest(action, reader)
	appCtx := &libApp.Context{}
	libApp.InitContextWithRequest(appCtx, req)
	appCtx.Set(SerlvetDoContext, ctx)

	err := app.Run(appCtx, nil)
	if err != nil {
		return nil, err
	}

	return appCtx, nil
}

func getApp(serlvet *stdSerlvet) SerlvetApp {
	serlvet.appRW.RLock()
	defer serlvet.appRW.RUnlock()
	return serlvet.app
}

func setApp(serlvet *stdSerlvet, app SerlvetApp) {
	serlvet.appRW.Lock()
	defer serlvet.appRW.Unlock()
	serlvet.app = app
}
