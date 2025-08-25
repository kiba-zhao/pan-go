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

var ErrPeerClientTransportNotFound = errors.New("peer.PeerClient Error: Transport Not Found")
var ErrPeerClientAccessDenied = errors.New("peer.PeerClient Error: Access Denied")
var ErrPeerClientPeerPurgeError = errors.New("peer.PeerClient Error: Peer Purge Error")

type PeerRequestName = app.RequestName
type PeerRouter = app.AppHandleGroup
type PeerContext = app.AppContext
type PeerNext = app.Next
type PeerRequest = *app.Request
type PeerResponse = *app.Response
type PeerHeaderItem = app.HeaderItem

type PeerTransport interface {
	RoundTrip(context.Context, PeerID, io.Reader) (io.ReadCloser, error)
	CanReach(PeerID) bool
}

type PeerTransportPurgeable interface {
	Purge(PeerID) error
}

type PeerClient interface {
	Do(context.Context, PeerID, PeerRequest) (PeerResponse, error)
	Request(context.Context, PeerID, PeerRequestName, io.Reader, ...PeerHeaderItem) (PeerResponse, error)
	CanReach(PeerID) bool
	Purge(PeerID) error

	RegisterPeerTransport(transport PeerTransport)
	UnregisterPeerTransport(transport PeerTransport)
}

type stdPeerClient struct {
	logger log.Logger

	transports   []PeerTransport
	transportsRW sync.RWMutex
}

var _ = (PeerClient)((*stdPeerClient)(nil))

func (client *stdPeerClient) Do(ctx context.Context, peerId PeerID, request PeerRequest) (PeerResponse, error) {
	transports := client.PeerTransports()
	if len(transports) <= 0 {
		return nil, ErrPeerClientTransportNotFound
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
		if errors.Is(err, ErrPeerClientAccessDenied) {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if resReader == nil {
		return nil, ErrPeerClientTransportNotFound
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

func (client *stdPeerClient) Request(ctx context.Context, peerId PeerID, name PeerRequestName, body io.Reader, headerItems ...PeerHeaderItem) (PeerResponse, error) {
	request := app.NewRequest(name, body)
	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}
	return client.Do(ctx, peerId, request)
}

// CanReach checks if a peer is reachable
func (client *stdPeerClient) CanReach(peerId PeerID) bool {
	transports := client.PeerTransports()
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
func (client *stdPeerClient) Purge(peerId PeerID) error {
	transports := client.PeerTransports()
	if len(transports) <= 0 {
		return ErrPeerClientTransportNotFound
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
			client.logger.Error("PeerCluster", "Purge Error: "+err.Error())
		}
	}

	if hasErr {
		return ErrPeerClientPeerPurgeError
	}
	return nil
}

func (client *stdPeerClient) PeerTransports() []PeerTransport {
	client.transportsRW.RLock()
	defer client.transportsRW.RUnlock()
	return client.transports
}

func (client *stdPeerClient) RegisterPeerTransport(transport PeerTransport) {
	client.transportsRW.Lock()
	defer client.transportsRW.Unlock()
	client.transports = append(client.transports, transport)
}

func (client *stdPeerClient) UnregisterPeerTransport(transport PeerTransport) {
	client.transportsRW.Lock()
	defer client.transportsRW.Unlock()

	client.transports = slices.DeleteFunc(client.transports, func(t PeerTransport) bool {
		return t == transport
	})
}
