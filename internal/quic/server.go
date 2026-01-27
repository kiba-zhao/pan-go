// Define peer server for quic
package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net"
	"pan/internal/log"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicServerUnavailable = errors.New("quic.QuicServer Error: Unavailable")

type stdQuicServer struct {
	logger  log.Logger
	network QuicNetwork

	port        uint16
	addrs       []string
	certificate tls.Certificate

	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	provider *stdQuicProvider
}

func (qs *stdQuicServer) Setup(config QuicConfig) {
	qs.logger.Debug("QuicServer", "Setup")

	qs.reloadLock.Lock()
	defer qs.reloadLock.Unlock()

	changed := false
	port := config.Port()
	addrs := config.Addrs()
	certificate := config.Certificate()

	if qs.port != port {
		qs.port = port
		changed = true
	}

	if !slices.Equal(qs.addrs, addrs) {
		qs.addrs = addrs
		changed = true
	}

	if !bytes.Equal(qs.certificate.Certificate[0], certificate.Certificate[0]) {
		qs.certificate = certificate
		changed = true
	}

	if !changed || qs.reload {
		return
	}

	qs.reload = true
	qs.reloadChan <- struct{}{}

}

func (qs *stdQuicServer) ListenAndServe(ctx context.Context) error {
	qs.logger.Debug("QuicServer", "ListenAndServe begin")
	defer qs.logger.Debug("QuicServer", "ListenAndServe end")

	var err error
	var wg sync.WaitGroup
	var port uint16
	var certBytes []byte

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-qs.reloadChan:
			qs.reloadLock.Lock()
			qs.reload = false
			qs.reloadLock.Unlock()
		}

		provider := qs.provider
		if err != nil {
			provider.RevokeAllTransports()
			break
		}

		certificate := qs.certificate
		netAddrs := qs.addrs
		netPort := qs.port

		if len(netAddrs) <= 0 || len(certificate.Certificate) <= 0 || netPort == 0 {
			provider.RevokeAllTransports()
			continue
		}

		if port != netPort || !bytes.Equal(certBytes, certificate.Certificate[0]) {
			provider.RevokeAllTransports()
		} else {
			transportsSeq := provider.SeqForTransports()
			for addr, _ := range transportsSeq {
				if idx := slices.Index(netAddrs, addr); idx >= 0 {
					//remove reserved ip in ipAddrs
					netAddrs = slices.Delete(netAddrs, idx, idx+1)
					continue
				}
				// close outdated transports
				provider.RevokeTransport(addr)
			}
		}

		port = netPort
		certBytes = certificate.Certificate[0]
		wg.Wait()

		if len(netAddrs) <= 0 {
			continue
		}

		network := qs.network
		// attach new transports and listen for serve
		for _, netAddr := range netAddrs {

			transport, err := provider.NewTransport(netAddr, netPort)
			if err != nil {
				qs.logger.Error("QuicServer", "provider.NewTransport Error: "+err.Error())
				continue
			}

			tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
			quicConf := &quic.Config{}

			ln, lnErr := transport.Listen(tlsConf, quicConf)
			if lnErr != nil {
				qs.logger.Error("QuicServer", "quic.ListenAddr Error: "+lnErr.Error())
				continue
			} else {
				qs.logger.Info("QuicServer", "quic.ListenAddr Success: "+netAddr)
			}

			wg.Add(1)
			go func(ln *quic.Listener, trConn net.PacketConn) {
				defer wg.Done()
				defer trConn.Close()
				for {
					conn, err := ln.Accept(ctx)
					if err != nil && conn != nil {
						conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
					}
					if err == nil {
						_, err = network.Serve(conn, nil)
					}
					if errors.Is(err, quic.ErrServerClosed) || errors.Is(err, context.Canceled) {
						break
					}
					if ctxErr := ctx.Err(); ctxErr != nil {
						break
					}
					qs.logger.Error("QuicServer", "quic.Listener.Accept Error: "+err.Error())
				}
			}(ln, transport.Conn)
		}
	}
	return err
}
