package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"pan/lib/peer"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicClusterUnavailable = errors.New("quic.QuicCluster Error: Unavailable")
var ErrQuicClusterRouteNotFound = errors.New("quic.QuicCluster Error: Route Not Found")
var ErrQuicClusterInvalidConnection = errors.New("quic.QuicCluster Error: Invalid Connection")
var ErrQuicClusterInvalidPeerID = errors.New("quic.QuicCluster Error: Invalid Peer ID")
var ErrQuicClusterConnectionNotFound = errors.New("quic.QuicCluster Error: Connection Not Found")

const (
	InternalErrorCode     quic.ApplicationErrorCode = 0x1
	NotAvailableErrorCode quic.ApplicationErrorCode = 0x2
	AccessDeniedErrorCode quic.ApplicationErrorCode = 0x3
	InvalidPeerErrorCode  quic.ApplicationErrorCode = 0x4
)

type QuicCluster interface {
	peer.PeerNetwork
	peer.PeerNetworkPurgeable

	Lookup(peerId peer.PeerID) QuicConn
	Dial(context.Context, peer.PeerID) (QuicConn, error)

	Serve(quic.Connection, peer.PeerID) (QuicConn, error)

	Invite(peerId peer.PeerID) (QuicConn, error)
	Do(context.Context, QuicConn, io.Reader) (quic.Stream, error)
	Route(peerId peer.PeerID, addr string) error

	ServeUDPConn() *net.UDPConn
}

type stdQuicCluster struct {
	certificate   tls.Certificate
	certificateRW sync.RWMutex

	peerCluster   peer.PeerCluster
	peerClusterRW sync.RWMutex

	agent  *stdQuicAgent
	server *stdQuicServer

	networkRW sync.RWMutex
	connMgr   *stdQuicConnMgr
	routeMgr  *stdQuicRouteMgr
}

var _ = (QuicCluster)((*stdQuicCluster)(nil))

func (cluster *stdQuicCluster) CanReach(peerId peer.PeerID) bool {
	canReach := false
	connArr := cluster.connMgr.Search(peerId)
	if len(connArr) > 0 {
		for _, conn := range connArr {
			canReach = !conn.Closed()
			if canReach {
				break
			}
		}
	}

	if !canReach {
		if route := cluster.routeMgr.Search(peerId); route != nil {
			canReach = route.Available()
		}
	}
	return canReach
}

func (cluster *stdQuicCluster) RoundTrip(ctx context.Context, peerId peer.PeerID, reader io.Reader) (io.ReadCloser, error) {
	conn := cluster.Lookup(peerId)
	if conn == nil {
		dialConn, err := cluster.Dial(ctx, peerId)
		if err == nil {
			conn = dialConn
		}
	}

	if conn == nil {
		inviteConn, err := cluster.Invite(peerId)
		if err != nil {
			return nil, err
		}
		conn = inviteConn
	}

	if conn == nil {
		return nil, ErrQuicClusterRouteNotFound
	}

	resReader, err := cluster.Do(ctx, conn, reader)

	if err != nil {
		appErr, ok := err.(*quic.ApplicationError)
		if ok && appErr.ErrorCode == AccessDeniedErrorCode {
			return nil, peer.ErrPeerClusterAccessDenied
		}
	}
	return resReader, err
}

func (cluster *stdQuicCluster) Purge(peerId peer.PeerID) error {
	cluster.networkRW.Lock()
	defer cluster.networkRW.Unlock()

	// remove route with peerId
	route := cluster.routeMgr.Search(peerId)
	if route != nil {
		cluster.routeMgr.Delete(route)
	}
	//

	// remove all conn with peerId
	cluster.connMgr.Clean(peerId)
	//

	return nil
}

func (cluster *stdQuicCluster) Lookup(peerId peer.PeerID) QuicConn {
	var conn QuicConn
	connArr := cluster.connMgr.Search(peerId)
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

func (cluster *stdQuicCluster) Dial(ctx context.Context, peerId peer.PeerID) (QuicConn, error) {
	cluster.networkRW.RLock()
	defer cluster.networkRW.RUnlock()

	route := cluster.routeMgr.Search(peerId)
	if route == nil || !route.Available() {
		return nil, ErrQuicClusterRouteNotFound
	}
	addrs := route.Addrs()
	if len(addrs) <= 0 {
		return nil, ErrQuicClusterRouteNotFound
	}

	var wg sync.WaitGroup
	var locker sync.Mutex
	var dialConn QuicConn
	dialCtx, dialCancel := context.WithCancelCause(ctx)
	defer dialCancel(nil)

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
			conn, err := dialAddr(dialCtx, cluster.serveTransport(), addr, cluster.Certificate())
			if err == nil {
				serveConn, err = cluster.Serve(conn, peerId)
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

func (cluster *stdQuicCluster) Certificate() tls.Certificate {
	cluster.certificateRW.RLock()
	defer cluster.certificateRW.RUnlock()
	return cluster.certificate
}

func (cluster *stdQuicCluster) SetCertificate(certificate tls.Certificate) {
	cluster.certificateRW.Lock()
	defer cluster.certificateRW.Unlock()
	cluster.certificate = certificate
}

func (cluster *stdQuicCluster) Serve(conn quic.Connection, peerId peer.PeerID) (QuicConn, error) {
	if conn == nil {
		return nil, ErrQuicClusterInvalidConnection
	}

	peerCluster := cluster.PeerCluster()
	agent := cluster.agent
	if peerCluster == nil || agent == nil {
		conn.CloseWithError(NotAvailableErrorCode, "")
		return nil, ErrQuicClusterUnavailable
	}

	serveConn, ok := conn.(QuicConn)
	if !ok {
		appErrCode := InvalidPeerErrorCode
		connPeerID, err := parsePeerID(conn)
		if err == nil && len(peerId) > 0 && !bytes.Equal(peerId, connPeerID) {
			err = ErrQuicClusterInvalidPeerID
		}
		if err == nil && len(peerId) <= 0 {
			err = peerCluster.Access(connPeerID)
			if err != nil {
				appErrCode = AccessDeniedErrorCode
			}
		}
		if err != nil {
			conn.CloseWithError(appErrCode, err.Error())
			return nil, err
		}
		serveConn, _ = cluster.connMgr.SelectOrStore(&stdQuicConn{Connection: conn, peerId: connPeerID, mgr: cluster.connMgr})
	}

	err := agent.Follow(serveConn)
	if err != nil {
		return nil, err
	}

	go serveQuicConn(serveConn, peerCluster)

	return serveConn, nil
}

func (cluster *stdQuicCluster) PeerCluster() peer.PeerCluster {
	cluster.peerClusterRW.RLock()
	defer cluster.peerClusterRW.RUnlock()
	return cluster.peerCluster
}

func (cluster *stdQuicCluster) SetPeerCluster(peerCluster peer.PeerCluster) {
	cluster.peerClusterRW.Lock()
	defer cluster.peerClusterRW.Unlock()
	cluster.peerCluster = peerCluster
}

func (cluster *stdQuicCluster) Invite(peerId peer.PeerID) (QuicConn, error) {

	connArr := cluster.connMgr.Search(peerId)
	if len(connArr) <= 0 {
		return nil, ErrQuicClusterConnectionNotFound
	}

	var ctrlConn QuicConn
	for _, conn := range connArr {
		if !conn.Closed() {
			ctrlConn = conn
			break
		}
	}

	if ctrlConn == nil {
		return nil, ErrQuicClusterConnectionNotFound
	}

	return cluster.agent.Invite(ctrlConn)
}

func (cluster *stdQuicCluster) Do(ctx context.Context, conn QuicConn, reader io.Reader) (quic.Stream, error) {
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

func (cluster *stdQuicCluster) Route(peerId peer.PeerID, addr string) error {

	serveConn, err := cluster.route(peerId, addr)
	if err != nil || serveConn == nil {
		return err
	}

	agent := cluster.agent
	if agent != nil {
		err = cluster.agent.Greet(serveConn)
	}

	return err
}

func (cluster *stdQuicCluster) route(peerId peer.PeerID, addr string) (QuicConn, error) {

	peerCluster := cluster.PeerCluster()
	if peerCluster != nil {
		return nil, ErrQuicClusterUnavailable
	}

	if err := peerCluster.Access(peerId); err != nil {
		return nil, err
	}

	route := cluster.routeMgr.Search(peerId)
	if route == nil {
		nroute, _ := cluster.routeMgr.SearchOrStore(&stdQuicRoute{peerId: peerId})
		route = nroute
	}

	route.Lock()
	defer route.Unlock()

	if route.Contains(addr) {
		return nil, nil
	}

	var serveConn QuicConn
	conn, err := dialAddr(context.Background(), cluster.serveTransport(), addr, cluster.Certificate())
	if err == nil {
		serveConn, err = cluster.Serve(conn, peerId)
	}
	if err == nil {
		err = route.Store(addr)
		if errors.Is(err, ErrQuicPeerRouteDuplicateAddress) {
			err = nil
		}
	}

	return serveConn, err
}

func (cluster *stdQuicCluster) ServeUDPConn() *net.UDPConn {
	tr := cluster.serveTransport()
	if tr == nil {
		return nil
	}

	if udpConn, ok := tr.Conn.(*net.UDPConn); ok {
		return udpConn
	}
	return nil
}

func (cluster *stdQuicCluster) serveTransport() *quic.Transport {
	server := cluster.server
	if server == nil {
		return nil
	}
	return server.Transport()
}

func parsePeerID(conn quic.Connection) (peer.PeerID, error) {
	state := conn.ConnectionState()
	certificate := state.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

func dialAddr(ctx context.Context, tr *quic.Transport, addr string, certificate tls.Certificate) (quic.Connection, error) {
	if tr == nil {
		return nil, ErrQuicClusterUnavailable
	}
	remoteAddr, remoteAddrErr := net.ResolveUDPAddr("udp", addr)
	if remoteAddrErr != nil {
		return nil, remoteAddrErr
	}

	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}

	var conn quic.Connection
	var err error
	for i := 0; i < 3; i++ {
		conn, err = tr.Dial(ctx, remoteAddr, tlsConf, quicConf)
		if err == nil {
			break
		}
	}
	return conn, err
}

func serveQuicConn(conn QuicConn, peerCluster peer.PeerCluster) error {
	defer conn.CloseWithError(quic.ApplicationErrorCode(0), "")

	var err error
	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			break
		}

		go peerCluster.Serve(stream, conn.PeerID())
	}

	return err
}
