package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"pan/lib/log"
	"pan/lib/peer"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicNetworkUnavailable = errors.New("quic.QuicNetwork Error: Unavailable")
var ErrQuicNetworkRouteNotFound = errors.New("quic.QuicNetwork Error: Route Not Found")
var ErrQuicNetworkConnectionNotFound = errors.New("quic.QuicNetwork Error: Connection Not Found")
var ErrQuicNetworkPeerGuardDeny = errors.New("quic.QuicNetwork Error: Peer Guard Deny")
var ErrQuicNetworkInvalidConnection = errors.New("quic.QuicNetwork Error: Invalid Connection")
var ErrQuicNetworkInvalidPeerID = errors.New("quic.QuicNetwork Error: Invalid Peer ID")
var ErrQuicNetworkServeTooMany = errors.New("quic.QuicNetwork Error: Serve Too Many")

const (
	InternalErrorCode     quic.ApplicationErrorCode = 0x1
	NotAvailableErrorCode quic.ApplicationErrorCode = 0x2
	AccessDeniedErrorCode quic.ApplicationErrorCode = 0x3
	InvalidPeerErrorCode  quic.ApplicationErrorCode = 0x4
	TooManyCode           quic.ApplicationErrorCode = 0x5
)

type QuicNetwork interface {
	Lookup(peerId peer.PeerID) QuicConn
	Dial(context.Context, peer.PeerID) (QuicConn, error)
	Do(context.Context, QuicConn, io.Reader) (quic.Stream, error)

	Invite(peerId peer.PeerID) (QuicConn, error)
	Route(peerId peer.PeerID, addr string) error

	Serve(quic.Connection, peer.PeerID) (QuicConn, error)
	SelectTransport(addr string) *QuicTransport
}

type stdQuicNetwork struct {
	logger log.Logger
	agent  *stdQuicAgent

	rw sync.RWMutex

	networkRW sync.RWMutex
	connMgr   *stdQuicConnMgr
	routeMgr  *stdQuicRouteMgr

	provider    *stdQuicProvider
	certificate tls.Certificate

	peerGuard   peer.PeerGuard
	peerGuardRW sync.RWMutex

	peerNetwork   peer.PeerNetwork
	peerNetworkRW sync.RWMutex

	dialThreshold uint16
	dialTimeout   time.Duration

	dialNum           uint16
	dialLocker        sync.Mutex
	dialPendingOnce   sync.Once
	dialPendingCh     chan struct{}
	dialPendingNum    uint16
	dialPendingLocker sync.Mutex
}

var _ = (peer.PeerTransport)((*stdQuicNetwork)(nil))

func (network *stdQuicNetwork) CanReach(peerId peer.PeerID) bool {
	canReach := false
	connArr := network.connMgr.Search(peerId)
	if len(connArr) > 0 {
		for _, conn := range connArr {
			canReach = !conn.Closed()
			if canReach {
				break
			}
		}
	}

	if !canReach {
		if route := network.routeMgr.Search(peerId); route != nil {
			canReach = route.Available()
		}
	}
	return canReach
}

func (network *stdQuicNetwork) RoundTrip(ctx context.Context, peerId peer.PeerID, reader io.Reader) (io.ReadCloser, error) {
	conn := network.Lookup(peerId)
	if conn == nil {
		dialConn, err := network.Dial(ctx, peerId)
		if err == nil {
			conn = dialConn
		}
	}

	if conn == nil {
		inviteConn, err := network.Invite(peerId)
		if err != nil {
			return nil, err
		}
		conn = inviteConn
	}

	if conn == nil {
		return nil, ErrQuicNetworkRouteNotFound
	}

	resReader, err := network.Do(ctx, conn, reader)

	if err != nil {
		appErr, ok := err.(*quic.ApplicationError)
		if ok && appErr.ErrorCode == AccessDeniedErrorCode {
			return nil, peer.ErrPeerNetworkAccessDenied
		}
	}
	return resReader, err
}

var _ = (peer.PeerTransportPurgeable)((*stdQuicNetwork)(nil))

func (network *stdQuicNetwork) Purge(peerId peer.PeerID) error {
	network.networkRW.Lock()
	defer network.networkRW.Unlock()

	// remove route with peerId
	route := network.routeMgr.Search(peerId)
	if route != nil {
		network.routeMgr.Delete(route)
	}
	//

	// remove all conn with peerId
	network.connMgr.Clean(peerId)
	//

	return nil
}

var _ = (QuicNetwork)((*stdQuicNetwork)(nil))

func (network *stdQuicNetwork) Lookup(peerId peer.PeerID) QuicConn {
	var conn QuicConn
	connArr := network.connMgr.Search(peerId)
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

func (network *stdQuicNetwork) Do(ctx context.Context, conn QuicConn, reader io.Reader) (quic.Stream, error) {
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
		if qStream, ok := stream.(*stdQuicStream); ok {
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

func (network *stdQuicNetwork) Dial(ctx context.Context, peerId peer.PeerID) (QuicConn, error) {
	network.networkRW.RLock()
	defer network.networkRW.RUnlock()

	network.rw.RLock()
	certificate := network.certificate
	network.rw.RUnlock()
	if len(certificate.Certificate) <= 0 {
		return nil, ErrQuicNetworkUnavailable
	}

	route := network.routeMgr.Search(peerId)
	if route == nil || !route.Available() {
		return nil, ErrQuicNetworkRouteNotFound
	}
	addrs := route.Addrs()
	if len(addrs) <= 0 {
		return nil, ErrQuicNetworkRouteNotFound
	}

	var dialCtx context.Context
	network.rw.RLock()
	if network.dialTimeout > 0 {
		dialCtx, _ = context.WithTimeout(ctx, network.dialTimeout)
	} else {
		dialCtx = ctx
	}
	network.rw.RUnlock()

	err := network.requestDial(dialCtx)
	if err != nil {
		return nil, err
	}
	defer network.releaseDial()

	var dialConn QuicConn
	for _, addr := range addrs {
		conn, connErr := dialAddr(dialCtx, network.provider, addr, certificate)
		err = connErr
		if err == nil {
			dialConn, err = network.Serve(conn, peerId)
			if err == nil || errors.Is(err, dialCtx.Err()) {
				break
			}
			if errors.Is(err, ErrQuicNetworkInvalidPeerID) || errors.Is(err, ErrQuicNetworkPeerGuardDeny) {
				route.Delete(addr)
			}
		} else if errors.Is(err, dialCtx.Err()) {
			break
		} else {
			route.Delete(addr)
		}
	}

	if dialConn != nil {
		return dialConn, nil
	}
	return dialConn, err
}

func (network *stdQuicNetwork) dialPendingChan() chan struct{} {
	network.dialPendingOnce.Do(func() {
		network.dialPendingCh = make(chan struct{})
	})
	return network.dialPendingCh
}

func (network *stdQuicNetwork) requestDial(ctx context.Context) error {

	network.dialLocker.Lock()
	network.dialNum += 1
	dialNum := network.dialNum

	network.rw.RLock()
	dialThreshold := network.dialThreshold
	network.rw.RUnlock()

	if dialNum > dialThreshold {
		network.dialPendingLocker.Lock()
		network.dialPendingNum += 1
		network.dialPendingLocker.Unlock()
	}
	network.dialLocker.Unlock()

	var err error
	if dialNum > dialThreshold {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case network.dialPendingChan() <- struct{}{}:
		}
		network.dialPendingLocker.Lock()
		if network.dialPendingNum > 0 {
			network.dialPendingNum -= 1
		}
		network.dialPendingLocker.Unlock()
	} else {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	if err != nil {
		network.dialLocker.Lock()
		if network.dialNum > 0 {
			network.dialNum -= 1
		}
		network.dialLocker.Unlock()
	}
	return err
}

func (network *stdQuicNetwork) releaseDial() {
	network.dialPendingLocker.Lock()
	dialPendingNum := network.dialPendingNum
	network.dialPendingLocker.Unlock()

	if dialPendingNum > 0 {
		select {
		case <-time.After(time.Millisecond * 600):
		case <-network.dialPendingChan():
		}
	}

	network.dialLocker.Lock()
	if network.dialNum > 0 {
		network.dialNum -= 1
	}
	network.dialLocker.Unlock()
}

func (network *stdQuicNetwork) Invite(peerId peer.PeerID) (QuicConn, error) {

	connArr := network.connMgr.Search(peerId)
	if len(connArr) <= 0 {
		return nil, ErrQuicNetworkConnectionNotFound
	}

	var ctrlConn QuicConn
	for _, conn := range connArr {
		if !conn.Closed() {
			ctrlConn = conn
			break
		}
	}

	if ctrlConn == nil {
		return nil, ErrQuicNetworkConnectionNotFound
	}

	return network.agent.Invite(ctrlConn)
}

func (network *stdQuicNetwork) Route(peerId peer.PeerID, addr string) error {

	serveConn, err := network.route(peerId, addr)
	if err != nil || serveConn == nil {
		return err
	}

	agent := network.agent
	if agent != nil {
		err = network.agent.Greet(serveConn)
	}

	return err
}

func (network *stdQuicNetwork) route(peerId peer.PeerID, addr string) (QuicConn, error) {

	network.rw.RLock()
	certificate := network.certificate
	network.rw.RUnlock()

	peerGuard := network.PeerGuard()
	if len(certificate.Certificate) <= 0 || peerGuard == nil {
		return nil, ErrQuicNetworkUnavailable
	}

	if !peerGuard.Access(peerId) {
		return nil, ErrQuicNetworkPeerGuardDeny
	}

	route := network.routeMgr.Search(peerId)
	if route == nil {
		nroute, _ := network.routeMgr.SearchOrStore(&stdQuicRoute{peerId: peerId})
		route = nroute
	}

	route.Lock()
	defer route.Unlock()

	if route.Contains(addr) {
		return nil, nil
	}

	var serveConn QuicConn
	conn, err := dialAddr(context.Background(), network.provider, addr, certificate)
	if err == nil {
		serveConn, err = network.Serve(conn, peerId)
	}
	if err == nil {
		err = route.Store(addr)
		if errors.Is(err, ErrQuicPeerRouteDuplicateAddress) {
			err = nil
		}
	}

	return serveConn, err
}

func (network *stdQuicNetwork) PeerGuard() peer.PeerGuard {
	network.peerGuardRW.RLock()
	defer network.peerGuardRW.RUnlock()
	return network.peerGuard
}

func (network *stdQuicNetwork) SetupPeerGuard(peerGuard peer.PeerGuard) {
	network.peerGuardRW.Lock()
	defer network.peerGuardRW.Unlock()
	network.peerGuard = peerGuard
}

func (network *stdQuicNetwork) Serve(conn quic.Connection, peerId peer.PeerID) (QuicConn, error) {
	peerGuard := network.PeerGuard()
	peerNetwork := network.PeerNetwork()
	if peerGuard == nil || peerNetwork == nil {
		return nil, ErrQuicNetworkUnavailable
	}

	if conn == nil {
		return nil, ErrQuicNetworkInvalidConnection
	}

	serveConn, ok := conn.(QuicConn)
	if !ok {
		appErrCode := InvalidPeerErrorCode
		connPeerID, err := parsePeerID(conn)
		if err == nil && len(peerId) > 0 && !bytes.Equal(peerId, connPeerID) {
			err = ErrQuicNetworkInvalidPeerID
		}
		if err == nil && len(peerId) <= 0 {
			if !peerGuard.Access(connPeerID) {
				appErrCode = AccessDeniedErrorCode
				err = ErrQuicNetworkPeerGuardDeny
			}
		}
		if err != nil {
			conn.CloseWithError(appErrCode, err.Error())
			return nil, err
		}
		serveConn, _ = network.connMgr.SelectOrStore(&stdQuicConn{Connection: conn, peerId: connPeerID, mgr: network.connMgr})
	}

	err := network.agent.Follow(serveConn)
	if err != nil {
		return nil, err
	}

	go serveQuicConn(serveConn, peerNetwork)

	return serveConn, nil
}

func (network *stdQuicNetwork) PeerNetwork() peer.PeerNetwork {
	network.peerNetworkRW.RLock()
	defer network.peerNetworkRW.RUnlock()
	return network.peerNetwork
}

func (network *stdQuicNetwork) SetupPeerNetwork(peerNetwork peer.PeerNetwork) {
	network.peerNetworkRW.Lock()
	defer network.peerNetworkRW.Unlock()
	network.peerNetwork = peerNetwork
}

func (network *stdQuicNetwork) SelectTransport(addr string) *QuicTransport {
	provider := network.provider
	if provider == nil {
		return nil
	}

	transport, ok := provider.Select(addr)
	if !ok {
		return nil
	}
	return transport
}

func (network *stdQuicNetwork) Setup(config QuicConfig) error {
	network.rw.Lock()
	defer network.rw.Unlock()

	certificate := config.Certificate()
	dialTimeout := config.DialTimeout()
	dialThreshold := config.DialThreshold()

	if !bytes.Equal(network.certificate.Certificate[0], certificate.Certificate[0]) {
		network.certificate = certificate
	}

	if network.dialTimeout != dialTimeout {
		network.dialTimeout = dialTimeout
	}

	if network.dialThreshold != dialThreshold {
		if network.dialThreshold < dialThreshold {
			growthNum := dialThreshold - network.dialThreshold
			network.dialPendingLocker.Lock()
			if growthNum > network.dialPendingNum {
				growthNum = network.dialPendingNum
			}
		growth_loop:
			for i := uint16(0); i < growthNum; i++ {
				select {
				case <-time.After(time.Millisecond * 600):
					break growth_loop
				case <-network.dialPendingChan():
				}
			}
			network.dialPendingLocker.Unlock()
		}
		network.dialThreshold = dialThreshold
	}

	return nil
}

func serveQuicConn(conn QuicConn, peerNetwork peer.PeerNetwork) error {
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go peerNetwork.Serve(stream, conn.PeerID())
	}

	return err
}

func dialAddr(ctx context.Context, provider *stdQuicProvider, addr string, certificate tls.Certificate) (quic.Connection, error) {

	remoteAddr, remoteAddrErr := net.ResolveUDPAddr("udp", addr)
	if remoteAddrErr != nil {
		return nil, remoteAddrErr
	}

	var conn quic.Connection
	var err error
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	transports := provider.SeqForTransportsWithAddrIP(remoteAddr.IP)

transports_loop:
	for _, transport := range transports {
		var timer <-chan time.Time

		for i := 0; i < 3; i++ {
			if timer != nil {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-timer:
				}
			}
			conn, err = transport.Dial(ctx, remoteAddr, tlsConf, quicConf)
			if err == nil || errors.Is(err, ctx.Err()) {
				break transports_loop
			}
			timer = time.After(300 * time.Millisecond)
		}

	}

	return conn, err
}

func parsePeerID(conn quic.Connection) (peer.PeerID, error) {
	state := conn.ConnectionState()
	certificate := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}
