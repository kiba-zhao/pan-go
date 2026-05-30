package net

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"iter"
	"net"
	"pan/internal/log"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicNetworkUnavailable = errors.New("net.QuicNetwork Error: Unavailable")
var ErrQuicNetworkPeerNotFound = errors.New("net.QuicNetwork Error: Peer Not Found")
var ErrQuicNetworkQuicServerUnavailable = errors.New("net.QuicNetwork Error: Quic Server Unavailable")
var ErrQuicPeerRouteDuplicateAddress = errors.New("net.QuicNetwork Error: Peer Route Duplicate Address")

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
	logger     log.Logger
	quicServer *stdQuicServer

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

func (network *stdQuicNetwork) SetupGuides(guides []QuicGuide) {
	network.guidesRW.Lock()
	defer network.guidesRW.Unlock()
	network.guides = guides
}

func (network *stdQuicNetwork) Setup(config QuicConfig) {
	network.logger.Debug("net.QuicNetwork", "Setup begin")
	defer network.logger.Debug("net.QuicNetwork", "Setup end")
	network.rw.Lock()
	defer network.rw.Unlock()

	network.certificate = config.Certificate()
}

func (network *stdQuicNetwork) RoundTrip(ctx context.Context, peerId PeerID) (PeerStream, error) {
	var stream PeerStream
	var err error
	var connSet *stdQuicConnSet

	network.connRW.RLock()
	connIdx, connOK := slices.BinarySearchFunc(network.connSets, peerId, comparePeerIDForQuicConnSet)
	if connOK {
		connSet = network.connSets[connIdx]
	}
	network.connRW.RUnlock()

	if connOK {
		connSet.locker.Lock()
		connArr := connSet.conns
		if len(connArr) > 0 {
			offset := 0
			for idx, conn := range connArr {
				stream, err = conn.OpenStream(ctx)
				if err == nil || errors.Is(err, ctx.Err()) {
					break
				}
				offset = idx + 1
			}

			if offset > len(connArr) {
				connSet.conns = nil

				network.connRW.Lock()
				network.connSets = slices.Delete(network.connSets, connIdx, connIdx+1)
				network.connRW.Unlock()
			} else if offset > 0 {
				connSet.conns = slices.Clone(connArr[offset:])
			}
		}
		connSet.locker.Unlock()
	}

	if err == nil {
		return stream, err
	}

	conn, err := network.Connect(ctx, peerId)
	if err == nil {
		err = network.Reuse(conn)
	}

	return conn.OpenStream(ctx)
}

func (network *stdQuicNetwork) Connect(ctx context.Context, peerId PeerID) (PeerConn, error) {
	conn, err := network.connectWithRoute(ctx, peerId)
	if err == nil || errors.Is(err, ctx.Err()) {
		return conn, err
	}
	return network.connectWithGuide(ctx, peerId)
}

func (network *stdQuicNetwork) Reuse(conn PeerConn) error {

	quicConn, err := parseQuicConn(conn)
	if err != nil {
		return err
	}

	peerId, err := extractPeerIDFromQuicConn(quicConn)
	if err != nil {
		return err
	}

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
	quicConnSet.conns = append(quicConnSet.conns, conn.(*stdQuicConn))

	return nil
}

func (network *stdQuicNetwork) Close(conn PeerConn) error {
	quicConn, err := parseQuicConn(conn)
	if err != nil {
		return err
	}
	return quicConn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
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
		conn, err = network.ConnectAddr(ctx, addr)
		if err == nil || errors.Is(err, ctx.Err()) {
			break
		}
		offset = addrIdx + 1
	}

	if offset > len(route.addrs) {
		route.addrs = nil

		network.routeRW.Lock()
		network.routes = slices.Delete(network.routes, idx, idx+1)
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
			conn, err = network.ConnectAddr(ctx, addr)
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

func (network *stdQuicNetwork) ConnectAddr(ctx context.Context, addr string) (*stdQuicConn, error) {
	remoteAddr, remoteAddrErr := net.ResolveUDPAddr("udp", addr)
	if remoteAddrErr != nil {
		return nil, remoteAddrErr
	}
	quicServer := network.quicServer
	if quicServer == nil {
		return nil, ErrQuicNetworkQuicServerUnavailable
	}

	network.rw.RLock()
	certificate := network.certificate
	network.rw.RUnlock()
	if len(certificate.Certificate) <= 0 {
		return nil, ErrQuicNetworkUnavailable
	}

	transports := quicServer.Transports()
	if len(transports) <= 0 {
		return nil, ErrQuicNetworkPeerNotFound
	}

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
		for i := 0; i < 3; i++ {
			if timer != nil {
				select {
				case <-ctx.Done():
					err = ctx.Err()
					break transports_loop
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

	if err != nil {
		return nil, err
	}
	return &stdQuicConn{Connection: conn}, nil
}

func (network *stdQuicNetwork) Route(quicConn *stdQuicConn) error {
	peerId, err := extractPeerIDFromQuicConn(quicConn)
	if err != nil {
		return err
	}

	addr := quicConn.RemoteAddr().String()

	network.routeRW.Lock()
	defer network.routeRW.Unlock()
	idx, ok := slices.BinarySearchFunc(network.routes, peerId, comparePeerIDForQuicRoute)
	if !ok {
		route := &stdQuicRoute{peerId: peerId}
		route.addrs = []string{addr}
		network.routes = slices.Insert(network.routes, idx, route)
		return nil
	}
	route := network.routes[idx]
	if slices.Contains(route.addrs, addr) {
		return ErrQuicPeerRouteDuplicateAddress
	}
	route.addrs = append(route.addrs, addr)
	return nil
}

func (network *stdQuicNetwork) HasRoute(peerId PeerID, addr string) bool {
	network.routeRW.RLock()
	defer network.routeRW.RUnlock()

	idx, ok := slices.BinarySearchFunc(network.routes, peerId, comparePeerIDForQuicRoute)
	if !ok {
		return false
	}
	route := network.routes[idx]
	return slices.Contains(route.addrs, addr)
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
