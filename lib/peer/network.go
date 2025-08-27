package peer

import (
	"context"
	"errors"
	"io"
	"pan/lib/app"
	"pan/lib/log"
	"slices"
	"sync"
)

var ErrPeerNetworkTransportNotFound = errors.New("peer.PeerNetwork Error: Transport Not Found")
var ErrPeerNetworkAccessDenied = errors.New("peer.PeerNetwork Error: Access Denied")
var ErrPeerNetworkPeerPurgeError = errors.New("peer.PeerNetwork Error: Peer Purge Error")
var ErrPeerNetworkNotFound = errors.New("peer.PeerNetwork Error: Not Found")
var ErrPeerNetworkInvalidApp = errors.New("peer.PeerNetwork Error: Invalid App")

type PeerRequestName = app.RequestName
type PeerRouter = app.AppHandleGroup
type PeerContext = app.AppContext
type PeerNext = app.Next
type PeerRequest = *app.Request
type PeerResponse = *app.Response
type PeerHeaderItem = app.HeaderItem

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

type PeerTransport interface {
	RoundTrip(context.Context, PeerID, io.Reader) (io.ReadCloser, error)
	CanReach(PeerID) bool
}

type PeerTransportPurgeable interface {
	Purge(PeerID) error
}

type PeerNetwork interface {
	Do(context.Context, PeerID, PeerRequest) (PeerResponse, error)
	Request(context.Context, PeerID, PeerRequestName, io.Reader, ...PeerHeaderItem) (PeerResponse, error)
	CanReach(PeerID) bool
	Purge(PeerID) error

	Serve(PeerStream, PeerID) error

	RegisterPeerTransport(transport PeerTransport)
	UnregisterPeerTransport(transport PeerTransport)
}

type stdPeerNetwork struct {
	logger log.Logger

	transports   []PeerTransport
	transportsRW sync.RWMutex

	peerApp   PeerApp
	peerAppRW sync.RWMutex
}

var _ = (PeerNetwork)((*stdPeerNetwork)(nil))

func (network *stdPeerNetwork) Do(ctx context.Context, peerId PeerID, request PeerRequest) (PeerResponse, error) {
	transports := network.PeerTransports()
	if len(transports) <= 0 {
		return nil, ErrPeerNetworkTransportNotFound
	}

	var resReader io.ReadCloser
	var err error
	reqReader := app.MarshalRequest(request)
	reader := &stdPeerReader{Reader: reqReader}
	for _, transport := range transports {

		if !transport.CanReach(peerId) {
			continue
		}

		resReader, err = transport.RoundTrip(ctx, peerId, reader)

		// break with success or request reader has been read
		if resReader != nil || reader.haveRead {
			break
		}

		// break with context is done
		if errors.Is(err, ctx.Err()) {
			break
		}

		// break access denied error
		if errors.Is(err, ErrPeerNetworkAccessDenied) {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if resReader == nil {
		return nil, ErrPeerNetworkTransportNotFound
	}

	response := &app.Response{}
	app.InitResponse(response)
	err = app.UnmarshalResponse(resReader, response)

	if err == nil && response.Code() != CodeOK {
		var content []byte
		content, err = io.ReadAll(response)
		if err == nil {
			err = &PeerError{code: response.Code(), err: string(content)}
		}
	}

	if err != nil {
		resReader.Close()
		return nil, err
	}

	return response, err
}

func (network *stdPeerNetwork) Request(ctx context.Context, peerId PeerID, name PeerRequestName, body io.Reader, headerItems ...PeerHeaderItem) (PeerResponse, error) {
	request := app.NewRequest(name, body)
	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}
	return network.Do(ctx, peerId, request)
}

// CanReach checks if a peer is reachable
func (network *stdPeerNetwork) CanReach(peerId PeerID) bool {
	transports := network.PeerTransports()
	if len(transports) <= 0 {
		return false
	}

	for _, transport := range transports {
		if transport.CanReach(peerId) {
			return true
		}
	}
	return false
}

// Purge purges a peer with the given peerId
func (network *stdPeerNetwork) Purge(peerId PeerID) error {
	transports := network.PeerTransports()
	if len(transports) <= 0 {
		return ErrPeerNetworkTransportNotFound
	}

	hasErr := false
	for _, transport := range transports {
		purgeableTransport, ok := transport.(PeerTransportPurgeable)
		if !ok {
			continue
		}
		err := purgeableTransport.Purge(peerId)
		if err != nil {
			hasErr = true
			network.logger.Error("PeerCluster", "Purge Error: "+err.Error())
		}
	}

	if hasErr {
		return ErrPeerNetworkPeerPurgeError
	}
	return nil
}

func (network *stdPeerNetwork) Serve(stream PeerStream, target PeerID) error {
	defer stream.Close()

	network.peerAppRW.RLock()
	peerApp := network.peerApp
	network.peerAppRW.RUnlock()

	var err error
	ctx := app.NewAppContext()
	if peerApp == nil {
		err = ErrPeerNetworkInvalidApp
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
		ctx.ThrowError(CodeNotFound, ErrPeerNetworkNotFound)
	}

	reader := app.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}

func (network *stdPeerNetwork) SetupApp(peerApp PeerApp) {
	network.peerAppRW.Lock()
	defer network.peerAppRW.Unlock()
	network.peerApp = peerApp
}

func (network *stdPeerNetwork) PeerTransports() []PeerTransport {
	network.transportsRW.RLock()
	defer network.transportsRW.RUnlock()
	return network.transports
}

func (network *stdPeerNetwork) RegisterPeerTransport(transport PeerTransport) {
	network.transportsRW.Lock()
	defer network.transportsRW.Unlock()
	network.transports = append(network.transports, transport)
}

func (network *stdPeerNetwork) UnregisterPeerTransport(transport PeerTransport) {
	network.transportsRW.Lock()
	defer network.transportsRW.Unlock()

	network.transports = slices.DeleteFunc(network.transports, func(t PeerTransport) bool {
		return t == transport
	})
}
