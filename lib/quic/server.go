// Define peer server for quic
package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicPeerServerUnavailable = errors.New("quic.QuicPeerServer Error: Unavailable")
var ErrQuicPeerSettingsUnavailable = errors.New("quic.QuicPeerServer Error: Peer Settings Unavailable")

type quicPeerServer struct {
	quicPeerModule QuicPeerModule
	locker         sync.RWMutex
	ln             *quic.Listener
	address        string
}

// Shutdown closes the quic listener and cleans up the quicPeerServer.
// Note that any new incoming connections will be closed immediately.
// Shutdown will return ErrQuicPeerServerUnavailable if the quicPeerServer
// is not available anymore.
func (qs *quicPeerServer) Shutdown() error {
	qs.locker.RLock()
	ln := qs.ln
	qs.locker.RUnlock()
	if ln == nil {
		return ErrQuicPeerServerUnavailable
	}

	qs.locker.Lock()
	qs.ln = nil
	qs.locker.Unlock()
	return ln.Close()
}

// ListenAndServe starts a QUIC server and listens for incoming connections.
// It authenticates clients using the server's TLS certificate and handles
// client connections using the associated quicPeerModule. The function will
// return an error if the server is unavailable, the peer settings are
// unavailable, or if there is an issue with accepting connections. It
// gracefully shuts down when the context is canceled or when the server
// encounters an internal error.

func (qs *quicPeerServer) ListenAndServe(ctx context.Context) error {

	if qs.quicPeerModule == nil {
		return ErrQuicPeerServerUnavailable
	}

	settings := qs.quicPeerModule.PeerSettings()
	if settings == nil || !settings.Available() {
		return ErrQuicPeerSettingsUnavailable
	}

	certificate := settings.Certificate()
	tlsConf := &tls.Config{ClientAuth: tls.RequireAnyClientCert, Certificates: []tls.Certificate{certificate}, InsecureSkipVerify: true, MinVersion: tls.VersionTLS13}
	quicConf := &quic.Config{}
	ln, err := quic.ListenAddr(qs.address, tlsConf, quicConf)
	if err != nil {
		return err
	}
	qs.locker.Lock()
	qs.ln = ln
	qs.locker.Unlock()
	defer qs.Shutdown()

	for {
		conn, err := ln.Accept(ctx)
		if err != nil && conn != nil {
			conn.CloseWithError(quic.ApplicationErrorCode(quic.InternalError), err.Error())
		}
		if err == nil {
			_, err = qs.quicPeerModule.Serve(conn, nil)
		}
		if errors.Is(err, quic.ErrServerClosed) || errors.Is(err, context.Canceled) {
			break
		}

	}

	return err
}
