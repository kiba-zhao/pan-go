package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"pan/app/bootstrap"
	"pan/app/config"
	"pan/app/peer"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicPeerModuleRouteConflict = errors.New("peer.PeerModule Error: Quic Peer Route Conflict")
var ErrPeerConflict = errors.New("peer.PeerModule Error: Peer Conflict")
var ErrPeerModuleUnknown = errors.New("peer.PeerModule Error: Unknown")
var ErrPeerSettingsUnavailable = errors.New("peer.PeerModule Error: Peer Settings Unavailable")
var ErrBroadcastDeliverExit = errors.New("quic.PeerBroadcast Error: Deliver Exit")

func parsePeerID(conn quic.Connection) (peer.PeerID, error) {
	state := conn.ConnectionState()
	certificate := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

type QuicPeerModule interface {
	PublicAddrs() []string
	PeerSettings() peer.PeerSettings
	Serve(quic.Connection) (peer.PeerNode, error)
	Do(context.Context, quic.Connection, io.Reader) (quic.Stream, error)
	Dial(context.Context, string) (quic.Connection, error)
	Route(peer.PeerID, string) error
}

type quicPeerModule struct {
	PeerModule        peer.PeerModule
	quicPeerBroadcast *quicPeerBroadcast

	publicAddrs     []string
	addrs           []string
	locker          sync.RWMutex
	reloadChan      chan struct{}
	reloadOnce      sync.Once
	reloadServe     bool
	reloadBroadcast bool

	mgr *quicPeerRouteMgr
	wg  sync.WaitGroup
}

func (qm *quicPeerModule) PeerSettings() peer.PeerSettings {
	if qm.PeerModule == nil {
		return nil
	}
	return qm.PeerModule.PeerSettings()
}

func (qm *quicPeerModule) PublicAddrs() []string {
	qm.locker.RLock()
	defer qm.locker.RUnlock()
	return qm.publicAddrs
}

func (qm *quicPeerModule) Addrs() []string {
	qm.locker.RLock()
	defer qm.locker.RUnlock()
	return qm.addrs
}

func (qm *quicPeerModule) ReloadChan() chan struct{} {
	qm.reloadOnce.Do(func() {
		qm.reloadChan = make(chan struct{}, 1)
	})
	return qm.reloadChan
}

func (qm *quicPeerModule) OnPeerSettingsUpdated(settings peer.PeerSettings) {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	if qm.reloadServe {
		return
	}
	qm.reloadServe = true
	qm.ReloadChan() <- struct{}{}
}

func (qm *quicPeerModule) OnConfigUpdated(settings config.AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	var reloadServe bool
	if reloadServe = !slices.Equal(qm.addrs, settings.PeerAddress); reloadServe {
		qm.addrs = settings.PeerAddress
	}

	if !qm.reloadServe && reloadServe {
		qm.reloadServe = reloadServe
	}

	var reloadBroadcast bool
	if reloadBroadcast = !slices.Equal(qm.publicAddrs, settings.PublicAddress); reloadBroadcast {
		qm.publicAddrs = settings.PublicAddress
	}

	if !qm.reloadBroadcast && reloadBroadcast {
		qm.reloadBroadcast = reloadBroadcast
	}

	if !qm.reloadServe && !qm.reloadBroadcast {
		return
	}

	qm.ReloadChan() <- struct{}{}
}

func (qm *quicPeerModule) Do(ctx context.Context, conn quic.Connection, reader io.Reader) (quic.Stream, error) {
	stream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	errChan := make(chan error)
	go func() {
		_, err = io.Copy(stream, reader)
		if err == nil {
			err = stream.Close()
		}

		errChan <- err
	}()

	select {
	case err = <-errChan:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return stream, err
}

func (qm *quicPeerModule) Dial(ctx context.Context, addr string) (quic.Connection, error) {
	if qm.PeerModule == nil {
		return nil, ErrPeerModuleUnknown
	}
	settings := qm.PeerSettings()
	if !settings.Available() {
		return nil, ErrPeerSettingsUnavailable
	}
	certificate := settings.Certificate()
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	return quic.DialAddr(ctx, addr, tlsConf, quicConf)
}

func (qm *quicPeerModule) createNode(conn quic.Connection) (*quicPeerNode, error) {

	peerId, err := parsePeerID(conn)
	if err != nil {
		return nil, err
	}

	peerModule := qm.PeerModule

	qnode := &quicPeerNode{
		quicPeerModule: qm,
		conn:           conn,
		peerId:         peerId,
		mgr:            peerModule.PeerManager(),
	}

	for {
		qnode.resourceId = peerModule.NewResourceID(peer.PeerTypeAlive)
		err = peerModule.Control(peer.PeerNode(qnode))
		if err == nil || err != peer.ErrPeerModuleControlConflict {
			break
		}
	}

	return qnode, err
}

func (qm *quicPeerModule) Route(peerId peer.PeerID, addr string) error {
	peerModule := qm.PeerModule

	routeId := make([]byte, 0)
	routeId = append(routeId, peerId...)
	routeId = append(routeId, []byte(addr)...)

	// check conflict
	route := qm.mgr.Search(routeId)
	if route != nil {
		return ErrQuicPeerModuleRouteConflict
	}
	//

	// try connect
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := qm.Dial(ctx, addr)
	if err == nil {
		connPeerId, parseErr := parsePeerID(conn)
		err = parseErr
		if err == nil && !bytes.Equal(connPeerId, peerId) {
			err = ErrPeerConflict
		}
		conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
	}
	if err != nil {
		return err
	}
	//

	route = &quicPeerRoute{
		quicPeerModule: qm,
		peerId:         peerId,
		address:        addr,
		routeId:        routeId,
	}

	for {
		route.resourceId = peerModule.NewResourceID(peer.PeerTypeReachable)
		err = peerModule.Control(peer.PeerNode(route))
		if err == nil || err != peer.ErrPeerModuleControlConflict {
			break
		}
	}

	// add to routes
	if err == nil {
		_, ok := qm.mgr.SearchOrStore(route)
		if ok {
			err = ErrQuicPeerModuleRouteConflict
		}
	}

	return err
}

func (qm *quicPeerModule) serve(peerNode *quicPeerNode) error {
	defer peerNode.Close()

	conn := peerNode.conn

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}
		peerNode.increaseStream()
		go qm.PeerModule.Serve(&quicPeerStream{Stream: stream, quicPeerNode: peerNode, isServe: true}, peerNode)
	}

	return err
}

func (qm *quicPeerModule) Serve(conn quic.Connection) (peer.PeerNode, error) {
	node, err := qm.createNode(conn)
	if err == nil {
		go qm.serve(node)
	}
	return node, err

}

func (qm *quicPeerModule) Components() []bootstrap.Component {
	return []bootstrap.Component{
		bootstrap.NewComponent(qm, bootstrap.ComponentNoneScope),
		bootstrap.NewComponent[QuicPeerModule](qm, bootstrap.ComponentExternalScope),
		bootstrap.NewComponent(qm.quicPeerBroadcast, bootstrap.ComponentNoneScope),
	}
}

func (qm *quicPeerModule) Ready(ctx context.Context) error {

	var servers []*quicPeerServer
	defer qm.shutdownForQuic(servers)

	var cancel context.CancelCauseFunc
	defer qm.shutdownForBroadcast(cancel)

	var err error
	closed := false

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-qm.ReloadChan():
		}

		qm.locker.Lock()
		reloadServe := qm.reloadServe
		qm.reloadServe = false

		reloadBroadcast := qm.reloadBroadcast
		qm.reloadBroadcast = false
		qm.locker.Unlock()

		if closed {
			break
		}

		if reloadServe {
			qm.shutdownForQuic(servers)
			servers = qm.serveForQuic(ctx)
		}

		if reloadBroadcast {
			qm.shutdownForBroadcast(cancel)
			causeCtx, causeCancel := context.WithCancelCause(ctx)
			cancel = causeCancel
			qm.quicPeerBroadcast.Ready(causeCtx)
		}

	}
	return err
}
func (qm *quicPeerModule) shutdownForQuic(servers []*quicPeerServer) {
	if len(servers) > 0 {
		for _, server := range servers {
			server.Shutdown()
		}
		qm.wg.Wait()
	}
}

func (qm *quicPeerModule) serveForQuic(ctx context.Context) []*quicPeerServer {
	servers := make([]*quicPeerServer, 0)
	addrs := qm.Addrs()
	for _, addr := range addrs {
		server := &quicPeerServer{
			address:        addr,
			quicPeerModule: qm,
		}

		servers = append(servers, server)
		qm.wg.Add(1)
		go func(s *quicPeerServer, c context.Context) {
			defer qm.wg.Done()
			_ = s.ListenAndServe(c)
			// TODO: write error into log
		}(server, ctx)
	}
	return servers
}

func (qm *quicPeerModule) shutdownForBroadcast(cancel context.CancelCauseFunc) {
	if cancel != nil {
		cancel(ErrBroadcastDeliverExit)
	}
}
