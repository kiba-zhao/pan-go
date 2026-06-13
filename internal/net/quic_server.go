package net

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"maps"
	"net"
	"pan/internal/log"
	"slices"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicServerInvalidQuicConn = errors.New("net.QuicServer Error: Invalid QuicConn")
var ErrQuicServerPeerGuardDenied = errors.New("net.QuicServer Error: PeerGuard Denied")

const QUIC_SERVER_LIMIT_MAX_SIZE = 256
const QUIC_SERVER_LIMIT_TIMEOUT = 6 * time.Second

type stdQuicServer struct {
	logger  log.Logger
	network *stdQuicNetwork
	guard   *stdPeerGuard

	peerServlet PeerServlet
	servletRW   sync.RWMutex

	port        uint16
	addrs       []string
	certificate tls.Certificate

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	transportConnList []*net.UDPConn
	transportMap      map[string]*quic.Transport
	transportsRW      sync.RWMutex

	limitIds    [][]byte
	limitLocker sync.Mutex
}

var _ = (PeerServer)((*stdQuicServer)(nil))

func (server *stdQuicServer) setupPeerServlet(peerServlet PeerServlet) {
	server.servletRW.Lock()
	defer server.servletRW.Unlock()
	server.peerServlet = peerServlet
}

func (server *stdQuicServer) setup(config QuicConfig) {
	server.logger.Debug("net.QuicServer", "Setup begin")
	defer server.logger.Debug("net.QuicServer", "Setup end")

	server.reloadLock.Lock()
	defer server.reloadLock.Unlock()

	changed := false
	port := config.Port()
	addrs := config.Addrs()
	certificate := config.Certificate()

	if server.port != port {
		server.port = port
		changed = true
	}

	if !slices.Equal(server.addrs, addrs) {
		server.addrs = addrs
		changed = true
	}

	if !bytes.Equal(server.certificate.Certificate[0], certificate.Certificate[0]) {
		server.certificate = certificate
		changed = true
	}

	if !changed || server.reload {
		return
	}

	server.reload = true
	server.reloadChan <- struct{}{}
}

func (server *stdQuicServer) listenAndServe(ctx context.Context) error {
	server.logger.Debug("net.QuicServer", "ListenAndServe begin")
	defer server.logger.Debug("net.QuicServer", "ListenAndServe end")

	var err error
	var wg sync.WaitGroup
	var port uint16
	var certBytes []byte

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-server.reloadChan:
			server.reloadLock.Lock()
			server.reload = false
			server.reloadLock.Unlock()
		}

		if err != nil {
			server.revokeAllTransports()
			break
		}

		certificate := server.certificate
		netAddrs := server.addrs
		netPort := server.port

		if len(netAddrs) <= 0 || len(certificate.Certificate) <= 0 || netPort == 0 {
			server.revokeAllTransports()
			continue
		}

		if port != netPort || !bytes.Equal(certBytes, certificate.Certificate[0]) {
			server.revokeAllTransports()
		} else {
			server.transportsRW.RLock()
			seqForTransportAddrs := maps.Keys(server.transportMap)
			server.transportsRW.RUnlock()

			for addr := range seqForTransportAddrs {
				if idx := slices.Index(netAddrs, addr); idx >= 0 {
					//remove reserved ip in ipAddrs
					netAddrs = slices.Delete(netAddrs, idx, idx+1)
					continue
				}
				// close outdated transports
				server.revokeTransport(addr)
			}
		}

		port = netPort
		certBytes = certificate.Certificate[0]
		wg.Wait()

		if len(netAddrs) <= 0 {
			continue
		}

		var transportMap map[string]*quic.Transport
		var transportConnList []*net.UDPConn
		if len(server.transportMap) > 0 {
			server.transportsRW.RLock()
			transportMap = maps.Clone(server.transportMap)
			transportConnList = slices.Clone(server.transportConnList)
			server.transportsRW.RUnlock()
		} else {
			transportMap = make(map[string]*quic.Transport)
			transportConnList = make([]*net.UDPConn, 0)
		}

		// attach new transports and listen for serve
		for _, netAddr := range netAddrs {

			_, transport, err := newQuicTransport(netAddr, netPort)
			if err != nil {
				server.logger.Error("net.QuicServer", "provider.NewTransport Error: "+err.Error())
				continue
			}

			tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
			quicConf := &quic.Config{}

			ln, lnErr := transport.Listen(tlsConf, quicConf)
			if lnErr != nil {
				server.logger.Error("net.QuicServer", "quic.ListenAddr Error: "+lnErr.Error())
				continue
			} else {
				server.logger.Info("net.QuicServer", "quic.ListenAddr Success: "+netAddr)
			}

			transportMap[netAddr] = transport
			transportConnList = append(transportConnList, transport.Conn.(*net.UDPConn))
			wg.Add(1)
			go func(ln *quic.Listener, trConn net.PacketConn) {
				defer wg.Done()
				defer trConn.Close()
				for {
					var peerId PeerID
					conn, err := ln.Accept(ctx)
					if err == nil {
						peerId, err = extractPeerIDFromQuicConnection(conn)
					}
					if err != nil && conn != nil {
						conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
					}
					if err == nil {
						quicConn := newQuicConn(conn, peerId)
						err = server.Serve(ctx, quicConn)
						if err != nil {
							conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
							continue
						}
					}
					if errors.Is(err, quic.ErrServerClosed) || errors.Is(err, context.Canceled) {
						break
					}
					if ctxErr := ctx.Err(); ctxErr != nil {
						break
					}
					server.logger.Error("net.QuicServer", "quic.Listener.Accept Error: "+err.Error())
				}
			}(ln, transport.Conn)
		}

		server.transportsRW.Lock()
		server.transportMap = transportMap
		server.transportConnList = transportConnList
		server.transportsRW.Unlock()
	}

	return err
}

func (server *stdQuicServer) Serve(ctx context.Context, conn PeerConn) error {
	quicConn, ok := conn.(*stdQuicConn)
	if !ok {
		return ErrQuicServerInvalidQuicConn
	}

	passport := server.guard.check(quicConn.PeerID())
	if passport == PeerGuardDeny {
		return ErrQuicServerPeerGuardDenied
	}

	var limitId []byte
	if passport == PeerGuardAllow {
		server.network.route(quicConn)
	} else {
		limitId = server.restrict(quicConn)
		if limitId == nil {
			return ErrQuicServerPeerGuardDenied
		}
	}

	go server.serveQuicConn(ctx, quicConn, limitId)
	return nil
}

func (server *stdQuicServer) serveQuicConn(ctx context.Context, conn *stdQuicConn, limitId []byte) error {
	defer conn.Close()

	var err error
	var ctx_ context.Context
	var cancel context.CancelFunc
	var passport uint8
	if len(limitId) > 0 {
		passport = PeerGuardDefault
		ctx_, cancel = context.WithTimeout(ctx, QUIC_SERVER_LIMIT_TIMEOUT)
		defer cancel()
	} else {
		passport = PeerGuardAllow
		ctx_ = ctx
	}

	for {
		stream, err := conn.AcceptStream(ctx_)
		if err == nil {
			err = server.serveStream(conn, stream, passport == PeerGuardAllow)
		}

		if passport == PeerGuardAllow {
			if err != nil {
				break
			}
			continue
		}

		if err != nil && ctx_.Err() != nil {
			break
		}

		passport = server.guard.check(conn.PeerID())
		if passport == PeerGuardDeny {
			err = ErrQuicServerPeerGuardDenied
			break
		}

		if passport == PeerGuardAllow {
			server.network.route(conn)
			server.release(limitId)

			ctx_ = ctx
		}
	}

	if passport != PeerGuardAllow {
		server.release(limitId)
	}
	return err
}

func (server *stdQuicServer) serveStream(conn *stdQuicConn, stream PeerStream, isConcurrently bool) error {

	server.servletRW.RLock()
	peerServlet := server.peerServlet
	server.servletRW.RUnlock()

	if !isConcurrently {
		return servePeerStream(peerServlet, conn, stream, conn.PeerID())
	}

	go func() {
		err := servePeerStream(peerServlet, conn, stream, conn.PeerID())
		if err != nil {
			conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), "")
		}
	}()
	return nil
}

func (server *stdQuicServer) revokeAllTransports() error {
	server.transportsRW.Lock()
	defer server.transportsRW.Unlock()
	if len(server.transportMap) <= 0 {
		return nil
	}

	transportMap := server.transportMap
	server.transportMap = nil
	server.transportConnList = nil

	errs := make([]error, 0)
	for _, transport := range transportMap {
		err := transport.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) <= 0 {
		return nil
	}
	return errors.Join(errs...)
}

func (server *stdQuicServer) revokeTransport(addr string) error {
	server.transportsRW.Lock()
	defer server.transportsRW.Unlock()
	if len(server.transportMap) <= 0 {
		return nil
	}
	transport, ok := server.transportMap[addr]
	if !ok {
		return nil
	}

	conn := transport.Conn.(*net.UDPConn)
	if ok {
		connIdx := slices.Index(server.transportConnList, conn)
		if connIdx >= 0 {
			server.transportConnList = slices.Delete(server.transportConnList, connIdx, connIdx+1)
		}
	}

	delete(server.transportMap, addr)
	return transport.Close()
}

func (server *stdQuicServer) getTransports() map[string]*quic.Transport {
	server.transportsRW.RLock()
	defer server.transportsRW.RUnlock()

	return maps.Clone(server.transportMap)
}

func (server *stdQuicServer) TransportConnList() []*net.UDPConn {
	server.transportsRW.RLock()
	defer server.transportsRW.RUnlock()

	return slices.Clone(server.transportConnList)
}

func (server *stdQuicServer) restrict(quicConn *stdQuicConn) []byte {
	var limitId []byte
	limitId = quicConn.RemoteAddr().(*net.UDPAddr).IP

	server.limitLocker.Lock()
	defer server.limitLocker.Unlock()
	if len(server.limitIds) >= QUIC_SERVER_LIMIT_MAX_SIZE {
		return nil
	}

	idx, ok := slices.BinarySearchFunc(server.limitIds, limitId, bytes.Compare)
	if ok {
		return nil
	}
	server.limitIds = slices.Insert(server.limitIds, idx, limitId)
	return limitId
}

func (server *stdQuicServer) release(limitId []byte) {
	server.limitLocker.Lock()
	defer server.limitLocker.Unlock()
	idx, ok := slices.BinarySearchFunc(server.limitIds, limitId, bytes.Compare)
	if ok {
		server.limitIds = slices.Delete(server.limitIds, idx, idx+1)
	}
}

func newQuicTransport(addr string, port uint16) (*net.IPNet, *quic.Transport, error) {
	var ipNet *net.IPNet
	var addrIP net.IP
	var err error
	if len(addr) > 0 {
		addrIP, ipNet, err = net.ParseCIDR(addr)
		if err != nil {
			return nil, nil, err
		}
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   addrIP,
		Port: int(port),
	})
	if err != nil {
		return nil, nil, err
	}

	transport := &quic.Transport{
		Conn: conn,
	}

	return ipNet, transport, err
}
