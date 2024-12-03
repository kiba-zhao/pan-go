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

	"github.com/quic-go/quic-go"
)

var ErrQuicPeerModuleRouteConflict = errors.New("peer.PeerModule Error: Quic Peer Route Conflict")
var ErrPeerConflict = errors.New("peer.PeerModule Error: Peer Conflict")
var ErrPeerModuleUnknown = errors.New("peer.PeerModule Error: Unknown")
var ErrPeerSettingsUnavailable = errors.New("peer.PeerModule Error: Peer Settings Unavailable")
var ErrBroadcastDeliverExit = errors.New("quic.PeerBroadcast Error: Deliver Exit")
var ErrPeerModuleRouteNotFound = errors.New("quic.PeerModule Error: Route Not Found")
var ErrPeerModuleInvalidPeerID = errors.New("quic.PeerModule Error: Invalid Peer ID")
var ErrPeerModuleInvalidConnection = errors.New("quic.PeerModule Error: Invalid Connection")
var ErrPeerModuleUnavailable = errors.New("quic.PeerModule Error: Unavailable")

func parsePeerID(conn quic.Connection) (peer.PeerID, error) {
	state := conn.ConnectionState()
	certificate := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

func dialAddr(ctx context.Context, addr string, quicPeerModule *quicPeerModule) (quic.Connection, error) {
	if quicPeerModule == nil {
		return nil, ErrPeerSettingsUnavailable
	}
	settings := quicPeerModule.PeerSettings()
	if !settings.Available() {
		return nil, ErrPeerSettingsUnavailable
	}
	certificate := settings.Certificate()
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}

	var conn quic.Connection
	var err error
	for i := 0; i < 3; i++ {
		conn, err = quic.DialAddr(ctx, addr, tlsConf, quicConf)
		if err == nil {
			break
		}
	}
	return conn, err
}

func serveQuicConn(conn QuicConn, peerModule peer.PeerModule, quicPeerBroadcast *quicPeerBroadcast) error {
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		// read flag
		flags := make([]byte, 1)
		_, err = stream.Read(flags)
		if err != nil {
			break
		}
		//

		if flags[0] == 0 {
			// serve with PeerModule
			go peerModule.Serve(stream, conn.PeerID())
			continue
		}

		// try reply greet
		if quicPeerBroadcast != nil {
			go quicPeerBroadcast.ReplyGreet(stream, conn)
		}
	}

	return err
}

type QuicPeerModule interface {
	peer.PeerNetwork
	PublicAddrs() []string
	PeerSettings() peer.PeerSettings
	Serve(quic.Connection, peer.PeerID) (QuicConn, error)
	Do(context.Context, QuicConn, io.Reader) (quic.Stream, error)
	Lookup(peer.PeerID) QuicConn
	Dial(context.Context, peer.PeerID) (QuicConn, error)
	Route(peer.PeerID, string) (QuicConn, error)
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

	networkRW sync.RWMutex
	connMgr   QuicConnMgr
	routeMgr  *quicRouteMgr
	wg        sync.WaitGroup
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

func (qm *quicPeerModule) CanReach(peerId peer.PeerID) bool {

	canReach := false
	connArr := qm.connMgr.Search(peerId)
	if len(connArr) > 0 {
		for _, conn := range connArr {
			canReach = conn.Available()
			if canReach {
				break
			}
		}
	}

	if !canReach {
		if route := qm.routeMgr.Search(peerId); route != nil {
			canReach = route.Available()
		}
	}
	return canReach
}

func (qm *quicPeerModule) RoundTrip(ctx context.Context, peerId peer.PeerID, reader io.Reader) (io.ReadCloser, error) {

	conn := qm.Lookup(peerId)
	if conn == nil {
		dialConn, err := qm.Dial(ctx, peerId)
		if err != nil {
			return nil, err
		}
		conn = dialConn
	}

	doReader := io.MultiReader(bytes.NewReader([]byte{0}), reader)
	return qm.Do(ctx, conn, doReader)
}

func (qm *quicPeerModule) Do(ctx context.Context, conn QuicConn, reader io.Reader) (quic.Stream, error) {
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
		if qStream, ok := stream.(*quicStream); ok {
			qStream.hangup = false
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

func (qm *quicPeerModule) Lookup(peerId peer.PeerID) QuicConn {
	var conn QuicConn
	connArr := qm.connMgr.Search(peerId)
	if len(connArr) > 0 {
		for _, connItem := range connArr {
			if !connItem.Available() {
				continue
			}
			conn = connItem
			break
		}
	}
	return conn
}

func (qm *quicPeerModule) Dial(ctx context.Context, peerId peer.PeerID) (QuicConn, error) {
	qm.networkRW.RLock()
	defer qm.networkRW.RUnlock()

	route := qm.routeMgr.Search(peerId)
	if route == nil || !route.Available() {
		return nil, ErrPeerModuleRouteNotFound
	}
	addrs := route.Addrs()
	if len(addrs) <= 0 {
		return nil, ErrPeerModuleRouteNotFound
	}

	var wg sync.WaitGroup
	var locker sync.Mutex
	var dialConn QuicConn
	dialCtx, dialCancel := context.WithCancelCause(ctx)

outer_loop:
	for _, addr := range addrs {
		select {
		case <-dialCtx.Done():
			break outer_loop
		default:
		}

		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			select {
			case <-dialCtx.Done():
				return
			default:
			}

			var serveConn QuicConn
			conn, err := dialAddr(dialCtx, addr, qm)
			if err == nil {
				serveConn, err = qm.Serve(conn, peerId)
			}
			if err != nil {
				route.Delete(addr)
				return
			}

			locker.Lock()
			defer locker.Unlock()
			if dialConn == nil {
				dialConn = serveConn
				dialCancel(err)
			}
		}(addr)
	}

	wg.Wait()
	return dialConn, dialCtx.Err()
}

func (qm *quicPeerModule) Route(peerId peer.PeerID, addr string) (QuicConn, error) {

	route := qm.routeMgr.Search(peerId)
	if route == nil {
		if err := qm.PeerModule.Access(peerId); err != nil {
			return nil, err
		}
		nroute, _ := qm.routeMgr.SearchOrStore(&quicRoute{peerId: peerId})
		route = nroute
	}

	route.Lock()
	defer route.Unlock()

	if route.Contains(addr) {
		return nil, nil
	}

	var serveConn QuicConn
	conn, err := dialAddr(context.Background(), addr, qm)
	if err == nil {
		serveConn, err = qm.Serve(conn, peerId)
	}
	if err != nil {
		return serveConn, err
	}

	return serveConn, route.Store(addr)
}

func (qm *quicPeerModule) Serve(conn quic.Connection, peerId peer.PeerID) (QuicConn, error) {
	if conn == nil {
		return nil, ErrPeerModuleInvalidConnection
	}

	if qm.PeerModule == nil {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "")
		return nil, ErrPeerModuleUnavailable
	}

	serveConn, ok := conn.(QuicConn)
	if !ok {
		connPeerID, err := parsePeerID(conn)
		if err == nil && len(peerId) > 0 && !bytes.Equal(peerId, connPeerID) {
			err = ErrPeerModuleInvalidPeerID
		}
		if err == nil && len(peerId) <= 0 {
			err = qm.PeerModule.Access(connPeerID)
		}
		if err != nil {
			conn.CloseWithError(quic.ApplicationErrorCode(0), "")
			return nil, err
		}
		serveConn, _ = qm.connMgr.SelectOrStore(&quicConn{Connection: conn, peerId: connPeerID})
	}

	go serveQuicConn(serveConn, qm.PeerModule, qm.quicPeerBroadcast)

	return serveConn, nil
}

func (qm *quicPeerModule) Purge(peerId peer.PeerID) error {

	qm.networkRW.Lock()
	defer qm.networkRW.Unlock()

	// remove route with peerId
	route := qm.routeMgr.Search(peerId)
	if route != nil {
		qm.routeMgr.Delete(route)
	}
	//

	// remove all conn with peerId
	qm.connMgr.Clean(peerId)
	//

	return nil
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
