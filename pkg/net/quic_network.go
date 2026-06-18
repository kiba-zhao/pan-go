package net

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"iter"
	"net"
	"pan/pkg/log"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicNetworkUnavailable = errors.New("net.QuicNetwork Error: Unavailable")
var ErrQuicNetworkPeerNotFound = errors.New("net.QuicNetwork Error: Peer Not Found")
var ErrQuicNetworkQuicServerUnavailable = errors.New("net.QuicNetwork Error: Quic Server Unavailable")
var ErrQuicPeerRouteDuplicateAddress = errors.New("net.QuicNetwork Error: Peer Route Duplicate Address")
var ErrQuicNetworkPeerConnConflict = errors.New("net.QuicNetwork Error: Peer Conn Conflict")
var ErrQuicNetworkPeerGuardDenied = errors.New("net.QuicNetwork Error: Peer Guard Forbidden")

const QUIC_NETWORK_MAX_ROUTES_SIZE = 16
const QUIC_NETWORK_MAX_RETRY_TIMES = 3
const QUIC_NETWORK_RETRY_INTERVAL = 300 * time.Millisecond
const QUIC_NETWORK_REUSE_CONN_THRESHOLD = 64

type QuicGuide interface {
	LookupQuicAddr(peerId PeerID) (iter.Seq[string], error)
}

type stdQuicConnSet struct {
	peerId PeerID

	conns  []*stdQuicConn
	locker sync.Mutex
}

type stdQuicRoute struct {
	peerId PeerID

	addrs []string
	rw    sync.RWMutex
}

type stdQuicNetwork struct {
	logger log.Logger
	server *stdQuicServer
	guard  *stdPeerGuard

	connSets []*stdQuicConnSet
	connRW   sync.RWMutex

	routes  []*stdQuicRoute
	routeRW sync.RWMutex

	guides   []QuicGuide
	guidesRW sync.RWMutex

	certificate tls.Certificate
	rw          sync.RWMutex
}

var _ = (PeerNetwork)((*stdQuicNetwork)(nil))

func (network *stdQuicNetwork) RoundTrip(ctx context.Context, peerId PeerID) (PeerStream, error) {
	if network.guard.check(peerId) != PeerGuardAllow {
		return nil, ErrQuicNetworkPeerGuardDenied
	}

	var stream PeerStream
	var connSet *stdQuicConnSet

	network.connRW.RLock()
	connIdx, connOK := slices.BinarySearchFunc(network.connSets, peerId, comparePeerIDForQuicConnSet)
	if connOK {
		connSet = network.connSets[connIdx]
	}
	network.connRW.RUnlock()

	var err error
	if connOK {
		connSet.locker.Lock()
		connArr := connSet.conns
		connArr_ := make([]*stdQuicConn, 0)
		totalNum := len(connArr)

		if totalNum <= 0 {
			goto UNLOCK_CONN_SET
		}

		for idx, conn := range connArr {
			if conn.isClosed() {
				continue
			}
			if !isIdleQuicConn(conn) {
				connArr_ = append(connArr_, conn)
				continue
			}
			stream, err = conn.OpenStream(ctx)
			if err == nil || ctx.Err() != nil {
				if idx <= 0 {
					connArr_ = connArr
				} else if idx < totalNum-1 {
					connArr_ = append(connArr_, connArr[idx:]...)
				}
				break
			}
		}

		if len(connArr_) >= totalNum {
			goto UNLOCK_CONN_SET
		}

		if len(connArr_) <= 0 {
			connSet.conns = nil

			network.connRW.Lock()
			connIdx, connOK = slices.BinarySearchFunc(network.connSets, peerId, comparePeerIDForQuicConnSet)
			if connOK && network.connSets[connIdx] == connSet {
				network.connSets = slices.Delete(network.connSets, connIdx, connIdx+1)
			}
			network.connRW.Unlock()
		} else {
			connSet.conns = connArr_
		}
		goto UNLOCK_CONN_SET

	UNLOCK_CONN_SET:
		connSet.locker.Unlock()
	}

	if err == nil {
		return stream, err
	}

	conn, err := network.connect(ctx, peerId)
	if err == nil {
		err = network.reuse(conn)
	}

	return conn.OpenStream(ctx)
}

func (network *stdQuicNetwork) Connect(ctx context.Context, peerId PeerID) (PeerConn, error) {
	if network.guard.check(peerId) != PeerGuardAllow {
		return nil, ErrQuicNetworkPeerGuardDenied
	}

	return network.connect(ctx, peerId)
}

func (network *stdQuicNetwork) Reuse(conn PeerConn) error {
	if network.guard.check(conn.PeerID()) != PeerGuardAllow {
		return ErrQuicNetworkPeerGuardDenied
	}
	quicConn, err := parseQuicConn(conn)
	if err != nil {
		return err
	}
	return network.reuse(quicConn)
}

func (network *stdQuicNetwork) reuse(quicConn *stdQuicConn) error {
	peerId := quicConn.PeerID()

	network.connRW.Lock()
	idx, ok := slices.BinarySearchFunc(network.connSets, peerId, comparePeerIDForQuicConnSet)
	var quicConnSet *stdQuicConnSet
	if ok {
		quicConnSet = network.connSets[idx]
	} else {
		quicConnSet = newQuicConnSet(peerId)
		network.connSets = slices.Insert(network.connSets, idx, quicConnSet)
	}
	network.connRW.Unlock()

	quicConnSet.locker.Lock()
	defer quicConnSet.locker.Unlock()
	if !slices.Contains(quicConnSet.conns, quicConn) {
		quicConnSet.conns = append(quicConnSet.conns, quicConn)
	}
	return nil
}

func (network *stdQuicNetwork) connect(ctx context.Context, peerId PeerID) (*stdQuicConn, error) {
	conn, err := network.connectWithRoute(ctx, peerId)
	if err != nil && !errors.Is(err, ctx.Err()) {
		conn, err = network.connectWithGuide(ctx, peerId)
	}
	return conn, err
}
func (network *stdQuicNetwork) connectWithRoute(ctx context.Context, peerId PeerID) (*stdQuicConn, error) {
	var route *stdQuicRoute
	network.routeRW.RLock()
	idx, ok := slices.BinarySearchFunc(network.routes, peerId, comparePeerIDForQuicRoute)
	if ok {
		route = network.routes[idx]
	}
	network.routeRW.RUnlock()

	if !ok {
		return nil, ErrQuicNetworkPeerNotFound
	}

	route.rw.Lock()
	defer route.rw.Unlock()

	if len(route.addrs) == 0 {
		return nil, ErrQuicNetworkPeerNotFound
	}

	var conn *stdQuicConn
	var err error
	offset := 0
	for addrIdx, addr := range route.addrs {
		conn, err = network.connectAddr(ctx, peerId, addr)
		if err == nil || errors.Is(err, ctx.Err()) {
			break
		}
		offset = addrIdx + 1
	}

	if offset > len(route.addrs) {
		route.addrs = nil

		network.routeRW.Lock()
		idx, ok = slices.BinarySearchFunc(network.routes, peerId, comparePeerIDForQuicRoute)
		if ok && network.routes[idx] == route {
			network.routes = slices.Delete(network.routes, idx, idx+1)
		}
		network.routeRW.Unlock()

		return nil, ErrQuicNetworkPeerNotFound
	} else if offset > 0 {
		route.addrs = slices.Clone(route.addrs[offset:])
	}

	return conn, err
}

func (network *stdQuicNetwork) connectWithGuide(ctx context.Context, peerId PeerID) (*stdQuicConn, error) {
	network.guidesRW.RLock()
	defer network.guidesRW.RUnlock()
	var conn *stdQuicConn
	var err error

guides_loop:
	for _, guide := range network.guides {
		addrSeq, err := guide.LookupQuicAddr(peerId)
		if err != nil {
			continue
		}
		for addr := range addrSeq {
			conn, err = network.connectAddr(ctx, peerId, addr)
			if err == nil || errors.Is(err, ctx.Err()) {
				break guides_loop
			}
		}
	}

	if err == nil && conn == nil {
		err = ErrQuicNetworkPeerNotFound
	}
	return conn, err
}

func (network *stdQuicNetwork) connectAddr(ctx context.Context, peerId PeerID, addr string) (*stdQuicConn, error) {
	remoteAddr, remoteAddrErr := net.ResolveUDPAddr("udp", addr)
	if remoteAddrErr != nil {
		return nil, remoteAddrErr
	}

	network.rw.RLock()
	certificate := network.certificate
	network.rw.RUnlock()
	if len(certificate.Certificate) <= 0 {
		return nil, ErrQuicNetworkUnavailable
	}

	transports := network.server.getTransports()
	if len(transports) <= 0 {
		return nil, ErrQuicNetworkPeerNotFound
	}

	var remotePeerId PeerID
	var conn quic.Connection
	var err error
	tlsConf := &tls.Config{Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	remoteIP := remoteAddr.IP

transports_loop:
	for localAddr, transport := range transports {

		_, localIPNet, localIPErr := net.ParseCIDR(localAddr)
		err = localIPErr
		if err != nil {
			continue
		}
		if (remoteIP.To4() == nil && localIPNet.IP.To4() != nil) || (remoteIP.To4() != nil && localIPNet.IP.To4() == nil) {
			continue
		}
		if !remoteIP.IsGlobalUnicast() && !localIPNet.Contains(remoteIP) {
			continue
		}
		if remoteIP.IsGlobalUnicast() && localIPNet.IP.IsLoopback() {
			continue
		}

		var timer <-chan time.Time
		for i := 0; i < QUIC_NETWORK_MAX_RETRY_TIMES; i++ {
			if timer != nil {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					break transports_loop
				case <-timer:
				}
			}
			conn, err = transport.Dial(ctx, remoteAddr, tlsConf, quicConf)
			if err == nil {
				remotePeerId, err = extractPeerIDFromQuicConnection(conn)
			}
			if err == nil || errors.Is(err, ctx.Err()) {
				break transports_loop
			}
			timer = time.After(QUIC_NETWORK_RETRY_INTERVAL)
		}
	}

	if err != nil {
		return nil, err
	}
	if !bytes.Equal(remotePeerId, peerId) {
		conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
		return nil, ErrQuicNetworkPeerConnConflict
	}

	quicConn := newQuicConn(conn, peerId)
	return quicConn, nil
}

func (network *stdQuicNetwork) route(quicConn *stdQuicConn) error {

	peerId := quicConn.PeerID()
	addr := quicConn.RemoteAddr().String()

	network.routeRW.Lock()
	idx, ok := slices.BinarySearchFunc(network.routes, peerId, comparePeerIDForQuicRoute)
	if !ok {
		defer network.routeRW.Unlock()
		route := &stdQuicRoute{peerId: peerId}
		route.addrs = []string{addr}
		network.routes = slices.Insert(network.routes, idx, route)

		return nil
	}
	route := network.routes[idx]
	network.routeRW.Unlock()

	route.rw.Lock()
	defer route.rw.Unlock()
	if !slices.Contains(route.addrs, addr) {
		if len(route.addrs) >= QUIC_NETWORK_MAX_ROUTES_SIZE {
			addrs := make([]string, 0, QUIC_NETWORK_MAX_ROUTES_SIZE)
			copy(addrs, route.addrs[1:])
			addrs[QUIC_NETWORK_MAX_ROUTES_SIZE-1] = addr
			route.addrs = addrs
		} else {
			route.addrs = append(route.addrs, addr)
		}
	}
	return nil
}

func (network *stdQuicNetwork) optimize(ctx context.Context) {
	network.connRW.RLock()
	connSets := slices.Clone(network.connSets)
	network.connRW.RUnlock()
	if len(connSets) <= 0 {
		return
	}

LOOP_CONN_SET:
	for _, connSet := range connSets {
		connSet.locker.Lock()

		connArr := make([]*stdQuicConn, 0)
		idleConnArr := make([]*stdQuicConn, 0)
		closedConnArr := make([]*stdQuicConn, 0)
		totalNum := len(connSet.conns)
		remainingNum := totalNum

		if totalNum <= 0 {
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		for _, conn := range connSet.conns {
			if ctx.Err() != nil {
				goto UNLOCK_CONN_SET_FOR_OPTIMIZE
			}
			if !conn.detect() {
				closedConnArr = append(closedConnArr, conn)
				remainingNum--
				continue
			}
			if isIdleQuicConn(conn) {
				idleConnArr = append(idleConnArr, conn)
				continue
			}
			connArr = append(connArr, conn)
		}

		if remainingNum == totalNum {
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		if remainingNum <= 0 {
			connSet.conns = nil

			network.connRW.Lock()
			connIdx, connOK := slices.BinarySearchFunc(network.connSets, connSet.peerId, comparePeerIDForQuicConnSet)
			if connOK && network.connSets[connIdx] == connSet {
				network.connSets = slices.Delete(network.connSets, connIdx, connIdx+1)
			}
			network.connRW.Unlock()
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		if len(connArr) <= 0 {
			if remainingNum >= QUIC_NETWORK_REUSE_CONN_THRESHOLD {
				connSet.conns = slices.Clone(idleConnArr[-1*remainingNum:])
			} else {
				connSet.conns = idleConnArr
			}
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		if len(connArr) >= QUIC_NETWORK_REUSE_CONN_THRESHOLD || len(idleConnArr) <= 0 {
			connSet.conns = connArr
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		if remainingNum < QUIC_NETWORK_REUSE_CONN_THRESHOLD {
			connSet.conns = append(idleConnArr, connArr...)
			goto UNLOCK_CONN_SET_FOR_OPTIMIZE
		}

		connSet.conns = make([]*stdQuicConn, 0, QUIC_NETWORK_REUSE_CONN_THRESHOLD)
		copy(connSet.conns, idleConnArr[len(connArr)-QUIC_NETWORK_REUSE_CONN_THRESHOLD:])
		copy(connSet.conns, connArr)
		goto UNLOCK_CONN_SET_FOR_OPTIMIZE

	UNLOCK_CONN_SET_FOR_OPTIMIZE:
		connSet.locker.Unlock()

		if len(closedConnArr) > 0 {
			for _, conn := range closedConnArr {
				conn.close()
			}
		}

		if ctx.Err() != nil {
			break LOOP_CONN_SET
		}

	}

}

func (network *stdQuicNetwork) setupGuides(guides []QuicGuide) {
	network.guidesRW.Lock()
	defer network.guidesRW.Unlock()
	network.guides = guides
}

func (network *stdQuicNetwork) setup(config QuicConfig) {
	network.logger.Debug("net.QuicNetwork", "Setup begin")
	defer network.logger.Debug("net.QuicNetwork", "Setup end")
	network.rw.Lock()
	defer network.rw.Unlock()

	network.certificate = config.Certificate()
}

func newQuicConnSet(peerId PeerID) *stdQuicConnSet {
	connSet := &stdQuicConnSet{peerId: peerId}
	connSet.conns = make([]*stdQuicConn, 0)
	return connSet
}

func comparePeerIDForQuicConnSet(connSet *stdQuicConnSet, peerId PeerID) int {
	return bytes.Compare(connSet.peerId, peerId)
}

func comparePeerIDForQuicRoute(route *stdQuicRoute, peerId PeerID) int {
	return bytes.Compare(route.peerId, peerId)
}
