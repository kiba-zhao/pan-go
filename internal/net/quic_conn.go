package net

import (
	"context"
	"crypto/x509"
	"errors"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicConnMaxStreamsExceeded = errors.New("net.QuicConn Error: Max Streams Exceeded")
var ErrQuicConnInvalidPeerConn = errors.New("net.QuicConn Error: Invalid PeerConn")

type stdQuicConn struct {
	quic.Connection
	state PeerConnState
	once  sync.Once
}

var _ = (PeerConn)((*stdQuicConn)(nil))

func (conn *stdQuicConn) ConnState() PeerConnState {
	conn.once.Do(func() {
		conn.state = &stdQuicConnState{conn: conn.Connection}
	})
	return conn.state
}

func (conn *stdQuicConn) OpenStream(ctx context.Context) (PeerStream, error) {
	stream, err := conn.Connection.OpenStream()
	if err != nil {
		return nil, err
	}
	return &stdQuicStream{Stream: stream}, nil
}

func (conn *stdQuicConn) AcceptStream(ctx context.Context) (PeerStream, error) {
	return conn.Connection.AcceptStream(ctx)
}

type stdQuicConnState struct {
	conn quic.Connection
}

func (state *stdQuicConnState) Session() ([]byte, error) {
	connectionState := state.conn.ConnectionState()
	certificate := connectionState.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

func parseQuicConn(conn PeerConn) (*stdQuicConn, error) {
	quicConn, ok := conn.(*stdQuicConn)
	if !ok {
		return nil, ErrQuicConnInvalidPeerConn
	}
	return quicConn, nil
}

func extractPeerIDFromQuicConn(quicConn *stdQuicConn) (PeerID, error) {
	return quicConn.ConnState().Session()
}
