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
