package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"pan/app/config"
	"pan/app/injection"
	"pan/app/peer"
	"pan/logger"
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
var ErrPeerModuleConnectionNotFound = errors.New("quic.PeerModule Error: Connection Not Found")

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

func serveQuicConn(conn QuicConn, peerModule peer.PeerModule) error {
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go peerModule.Serve(stream, conn.PeerID())
	}

	return err
}

type QuicPeerModule interface {
	peer.PeerNetwork
	PeerSettings() peer.PeerSettings
	Serve(quic.Connection, peer.PeerID) (QuicConn, error)
	Do(context.Context, QuicConn, io.Reader) (quic.Stream, error)
	Lookup(peer.PeerID) QuicConn
	Dial(context.Context, peer.PeerID) (QuicConn, error)
	Route(peer.PeerID, string, bool) error
	Invite(peer.PeerID) (QuicConn, error)
	Reload()
}

type quicPeerModule struct {
	PeerModule peer.PeerModule
	agent      *quicPeerAgent

	addrs       []string
	locker      sync.RWMutex
	reloadChan  chan struct{}
	reloadOnce  sync.Once
	reloadServe bool

	networkRW sync.RWMutex
	connMgr   *quicConnMgr
	routeMgr  *quicRouteMgr
	wg        sync.WaitGroup
}

func (qm *quicPeerModule) PeerSettings() peer.PeerSettings {
	if qm.PeerModule == nil {
		return nil
	}
	return qm.PeerModule.PeerSettings()
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

func (qm *quicPeerModule) Reload() {
	qm.locker.Lock()
	defer qm.locker.Unlock()

	reload(qm)
}

func (qm *quicPeerModule) OnPeerSettingsUpdated(settings peer.PeerSettings) {
	qm.Reload()
}

func (qm *quicPeerModule) OnConfigUpdated(settings config.AppSettings) {

	qm.locker.Lock()
	defer qm.locker.Unlock()

	if reloadServe := !slices.Equal(qm.addrs, settings.PeerAddress); reloadServe {
		qm.addrs = settings.PeerAddress
	}

	reload(qm)
}

func (qm *quicPeerModule) CanReach(peerId peer.PeerID) bool {

	canReach := false
	connArr := qm.connMgr.Search(peerId)
	if len(connArr) > 0 {
		for _, conn := range connArr {
			canReach = !conn.Closed()
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
		if err == nil {
			conn = dialConn
		}
	}

	if conn == nil {
		inviteConn, err := qm.Invite(peerId)
		if err != nil {
			return nil, err
		}
		conn = inviteConn
	}

	if conn == nil {
		return nil, ErrPeerModuleRouteNotFound
	}

	return qm.Do(ctx, conn, reader)
}

func (qm *quicPeerModule) Do(ctx context.Context, conn QuicConn, reader io.Reader) (quic.Stream, error) {
	stream, err := conn.OpenStream()
	if err != nil {
		return nil, err
	}

	errCh := make(chan error)
	defer close(errCh)

	go func(ch chan error) {
		_, err = io.Copy(stream, reader)
		if err == nil {
			err = stream.Close()
		}
		if qStream, ok := stream.(*quicStream); ok {
			qStream.hangup = false
		}

		ch <- err
	}(errCh)

	select {
	case <-ctx.Done():
		stream.Close()
		err = ctx.Err()
	case err = <-errCh:
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
	if dialConn != nil {
		return dialConn, nil
	}
	return dialConn, dialCtx.Err()
}

func (qm *quicPeerModule) Route(peerId peer.PeerID, addr string, needGreet bool) error {

	route := qm.routeMgr.Search(peerId)
	if route == nil {
		if err := qm.PeerModule.Access(peerId); err != nil {
			return err
		}
		nroute, _ := qm.routeMgr.SearchOrStore(&quicRoute{peerId: peerId})
		route = nroute
	}

	route.Lock()
	defer route.Unlock()

	if route.Contains(addr) {
		return nil
	}

	var serveConn QuicConn
	conn, err := dialAddr(context.Background(), addr, qm)
	if err == nil {
		serveConn, err = qm.Serve(conn, peerId)
	}
	if err == nil {
		err = route.Store(addr)
	}
	if err == nil && needGreet {
		err = qm.agent.Greet(serveConn)
	}

	return err
}

func (qm *quicPeerModule) Invite(peerId peer.PeerID) (QuicConn, error) {

	connArr := qm.connMgr.Search(peerId)
	if len(connArr) <= 0 {
		return nil, ErrPeerModuleConnectionNotFound
	}

	var ctrlConn QuicConn
	for _, conn := range connArr {
		if !conn.Closed() {
			ctrlConn = conn
			break
		}
	}

	if ctrlConn == nil {
		return nil, ErrPeerModuleConnectionNotFound
	}

	return qm.agent.Invite(ctrlConn)
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
		serveConn, _ = qm.connMgr.SelectOrStore(&quicConn{Connection: conn, peerId: connPeerID, mgr: qm.connMgr})
	}

	err := qm.agent.Follow(serveConn)
	if err != nil {
		return nil, err
	}

	go serveQuicConn(serveConn, qm.PeerModule)

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

func (qm *quicPeerModule) Components() []injection.Component {
	return []injection.Component{
		injection.NewComponent(qm, injection.ComponentNoneScope),
		injection.NewComponent[QuicPeerModule](qm, injection.ComponentExternalScope),
		injection.NewComponent(qm.agent, injection.ComponentNoneScope),
	}
}

func (qm *quicPeerModule) Ready(ctx context.Context) error {

	var servers []*quicPeerServer

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
		qm.reloadServe = false
		qm.locker.Unlock()

		qm.shutdownForQuic(servers)
		if closed {
			break
		}

		servers = qm.serveForQuic(ctx)

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
			err := s.ListenAndServe(c)
			if err != nil {
				logger.Default().Log(context.Background(), logger.LevelError, "app.quic.serveForQuic Error: %s", err.Error())
			}
		}(server, ctx)
	}
	return servers
}

func reload(qm *quicPeerModule) {
	if qm.reloadServe {
		return
	}
	qm.reloadServe = true
	qm.ReloadChan() <- struct{}{}
}
