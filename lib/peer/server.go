package peer

import (
	"errors"
	"io"
	"pan/lib/app"
	"pan/lib/log"
	"sync"
)

var ErrPeerServerInvalidApp = errors.New("peer.PeerServer Error: Invalid App")
var ErrPeerServerNotFound = errors.New("peer.PeerServer Error: Not Found")

const (
	CodeOK            = app.CodeOK
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeForbidden     = 403
)

var (
	ContextPeerID = []byte("PeerID")
)

type PeerID = []byte
type PeerApp = *app.App

type PeerStream interface {
	io.Reader
	io.Writer
	io.Closer
}

type PeerServer interface {
	Serve(PeerStream, PeerID) error
}

type stdPeerServer struct {
	logger log.Logger

	peerApp   PeerApp
	peerAppRW sync.RWMutex
}

var _ = (PeerServer)((*stdPeerServer)(nil))

func (server *stdPeerServer) Serve(stream PeerStream, target PeerID) error {
	defer stream.Close()

	peerApp := server.PeerApp()

	var err error
	ctx := app.NewAppContext()
	if peerApp == nil {
		err = ErrPeerServerInvalidApp
	} else {
		err = app.UnmarshalRequest(stream, ctx.Request())
	}

	if err == nil {
		ctx.Set(ContextPeerID, target)
		err = peerApp.Run(ctx, nil)
		defer ctx.Close()
	}

	if err != nil {
		ctx.ThrowError(CodeInternalError, err)
	}

	if ctx.Code() < 0 {
		ctx.ThrowError(CodeNotFound, ErrPeerServerNotFound)
	}

	reader := app.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}

func (server *stdPeerServer) Setup(peerApp PeerApp) {
	server.peerAppRW.Lock()
	defer server.peerAppRW.Unlock()
	server.peerApp = peerApp
}

func (server *stdPeerServer) PeerApp() PeerApp {
	server.peerAppRW.RLock()
	defer server.peerAppRW.RUnlock()
	return server.peerApp
}
