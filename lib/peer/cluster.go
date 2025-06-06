package peer

import (
	"context"
	"errors"
	"io"
	"pan/lib/app"
	"pan/lib/log"
	"sync"
)

const (
	CodeOK            = app.CodeOK
	CodeInternalError = 500
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeForbidden     = 403
)

var ErrPeerClusterNetworkNotFound = errors.New("peer.PeerCluster Error: Network Not Found")
var ErrPeerClusterPeerPurgeError = errors.New("peer.PeerCluster Error: Peer Purge Error")
var ErrPeerClusterInvalidApp = errors.New("peer.PeerCluster Error: Invalid App")
var ErrPeerClusterNotFound = errors.New("peer.PeerCluster Error: Not Found")
var ErrPeerClusterAccessDenied = errors.New("peer.PeerCluster Error: Access Denied")

type PeerID = []byte
type PeerApp = *app.App
type PeerRouter = app.AppHandleGroup
type PeerContext = app.AppContext
type PeerNext = app.Next
type PeerRequest = *app.Request
type PeerResponse = *app.Response
type PeerHeaderItem = app.HeaderItem
type PeerRequestName = app.RequestName

type PeerStream interface {
	io.Reader
	io.Writer
	io.Closer
}

var (
	ContextPeerID = []byte("PeerID")
)

type PeerNetwork interface {
	RoundTrip(context.Context, PeerID, io.Reader) (io.ReadCloser, error)
	CanReach(PeerID) bool
}

type PeerNetworkPurgeable interface {
	Purge(PeerID) error
}

// PeerGuard is a module that provides access control and security features within the p2p module
type PeerGuard interface {
	Enabled() bool
	Access(PeerID) error
}

type PeerCluster interface {
	Serve(PeerStream, PeerID) error

	RegisterPeerNetwork(PeerNetwork)
	UnregisterPeerNetwork(PeerNetwork)
	PeerNetworks() []PeerNetwork
	Do(context.Context, PeerID, PeerRequest) (PeerResponse, error)
	Request(context.Context, PeerID, PeerRequestName, io.Reader, ...PeerHeaderItem) (PeerResponse, error)
	CanReach(PeerID) bool
	Purge(PeerID) error

	RegisterPeerGuard(PeerGuard)
	UnregisterPeerGuard(PeerGuard)
	PeerGuards() []PeerGuard
	Access(PeerID) error
}

type stdPeerCluster struct {
	logger log.Logger

	peerApp   PeerApp
	peerAppRW sync.RWMutex

	networks   []PeerNetwork
	networksRW sync.RWMutex

	guards   []PeerGuard
	guardsRW sync.RWMutex
}

var _ = (PeerCluster)((*stdPeerCluster)(nil))

// func NewPeerCluster(logger log.Logger) *PeerCluster {
// 	cluser := &PeerCluster{}
// 	cluser.logger = logger
// 	return cluser
// }

func (cluster *stdPeerCluster) SetPeerApp(peerApp PeerApp) {
	cluster.peerAppRW.Lock()
	defer cluster.peerAppRW.Unlock()
	cluster.peerApp = peerApp
}

func (cluster *stdPeerCluster) PeerApp() PeerApp {
	cluster.peerAppRW.RLock()
	defer cluster.peerAppRW.RUnlock()
	return cluster.peerApp
}

func (cluster *stdPeerCluster) Serve(stream PeerStream, target PeerID) error {

	defer stream.Close()

	peerApp := cluster.PeerApp()

	var err error
	ctx := app.NewAppContext()
	if peerApp == nil {
		err = ErrPeerClusterInvalidApp
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
		ctx.ThrowError(CodeNotFound, ErrPeerClusterNotFound)
	}

	reader := app.MarshalResponse(&ctx.Response)
	_, resErr := io.Copy(stream, reader)

	if err == nil && resErr != nil {
		err = resErr
	}

	return err
}

func (cluster *stdPeerCluster) RegisterPeerNetwork(network PeerNetwork) {
	cluster.networksRW.Lock()
	defer cluster.networksRW.Unlock()
	cluster.networks = append(cluster.networks, network)
}

func (cluster *stdPeerCluster) UnregisterPeerNetwork(network PeerNetwork) {
	cluster.networksRW.Lock()
	defer cluster.networksRW.Unlock()
	for i, n := range cluster.networks {
		if n == network {
			cluster.networks = append(cluster.networks[:i], cluster.networks[i+1:]...)
			break
		}
	}
}

func (cluster *stdPeerCluster) PeerNetworks() []PeerNetwork {
	cluster.networksRW.RLock()
	defer cluster.networksRW.RUnlock()
	return cluster.networks
}

func (cluster *stdPeerCluster) Do(ctx context.Context, peerId PeerID, request PeerRequest) (PeerResponse, error) {
	networks := cluster.PeerNetworks()
	if len(networks) <= 0 {
		return nil, ErrPeerClusterNetworkNotFound
	}

	var resReader io.ReadCloser
	var err error
	reqReader := app.MarshalRequest(request)
	reader := &stdPeerReader{Reader: reqReader}
	for _, network := range networks {

		if !network.CanReach(peerId) {
			continue
		}

		resReader, err = network.RoundTrip(ctx, peerId, reader)

		// break with success or request reader has been read
		if resReader != nil || reader.haveRead {
			break
		}

		// break with context is done
		if errors.Is(err, ctx.Err()) {
			break
		}

		// break access denied error
		if errors.Is(err, ErrPeerClusterAccessDenied) {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	if resReader == nil {
		return nil, ErrPeerClusterNetworkNotFound
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

func (cluster *stdPeerCluster) Request(ctx context.Context, peerId PeerID, name PeerRequestName, body io.Reader, headerItems ...PeerHeaderItem) (PeerResponse, error) {
	request := app.NewRequest(name, body)
	if len(headerItems) > 0 {
		for _, headerItem := range headerItems {
			request.SetHeader(headerItem.Key, headerItem.Value)
		}
	}
	return cluster.Do(ctx, peerId, request)
}

// CanReach checks if a peer is reachable
func (cluster *stdPeerCluster) CanReach(peerId PeerID) bool {
	networks := cluster.PeerNetworks()
	if len(networks) <= 0 {
		return false
	}

	for _, network := range networks {
		if network.CanReach(peerId) {
			return true
		}
	}
	return false
}

// Purge purges a peer with the given peerId
func (cluster *stdPeerCluster) Purge(peerId PeerID) error {
	networks := cluster.PeerNetworks()
	if len(networks) <= 0 {
		return ErrPeerClusterNetworkNotFound
	}

	hasErr := false
	for _, network := range networks {
		purgeableNetwork, ok := network.(PeerNetworkPurgeable)
		if !ok {
			continue
		}
		err := purgeableNetwork.Purge(peerId)
		if err != nil {
			hasErr = true
			cluster.logger.Error("PeerCluster", "Purge Error: "+err.Error())
		}
	}

	if hasErr {
		return ErrPeerClusterPeerPurgeError
	}
	return nil
}

func (cluster *stdPeerCluster) RegisterPeerGuard(guard PeerGuard) {
	cluster.guardsRW.Lock()
	defer cluster.guardsRW.Unlock()
	cluster.guards = append(cluster.guards, guard)
}

func (cluster *stdPeerCluster) UnregisterPeerGuard(guard PeerGuard) {
	cluster.guardsRW.Lock()
	defer cluster.guardsRW.Unlock()
	for i, g := range cluster.guards {
		if g == guard {
			cluster.guards = append(cluster.guards[:i], cluster.guards[i+1:]...)
			break
		}
	}
}

func (cluster *stdPeerCluster) PeerGuards() []PeerGuard {
	cluster.guardsRW.RLock()
	defer cluster.guardsRW.RUnlock()
	return cluster.guards
}

func (cluster *stdPeerCluster) Access(peerId PeerID) error {

	var err error
	guards := cluster.PeerGuards()
	if len(guards) <= 0 {
		return err
	}

	for _, guard := range guards {
		if !guard.Enabled() {
			continue
		}
		err = guard.Access(peerId)
		if err != nil {
			break
		}
	}
	return err
}
