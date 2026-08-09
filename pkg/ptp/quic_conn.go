package ptp

import (
	"context"
	"crypto/x509"
	"errors"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

var ErrQuicConnMaxStreamsExceeded = errors.New("net.QuicConn Error: Max Streams Exceeded")
var ErrQuicConnInvalidPeerConn = errors.New("net.QuicConn Error: Invalid PeerConn")
var ErrQuicConnDisconnected = errors.New("net.QuicConn Error: Disconnected")

const QUIC_CONN_DETECT_TIMEOUT = 100 * time.Millisecond // 100ms

type stdQuicConn struct {
	quic.Connection
	once sync.Once

	peerId PeerID

	closed bool
	rw     sync.RWMutex

	detectLocker sync.Mutex

	openStreamCount int
	openStreamRW    sync.RWMutex
}

func newQuicConn(conn quic.Connection, peerId PeerID) *stdQuicConn {
	quicConn := &stdQuicConn{Connection: conn}
	quicConn.peerId = peerId

	return quicConn
}

var _ = (PeerConn)((*stdQuicConn)(nil))

func (conn *stdQuicConn) PeerID() PeerID {
	return conn.peerId
}

func (conn *stdQuicConn) OpenStream(ctx context.Context) (PeerStream, error) {
	stream, err := conn.Connection.OpenStream()
	if err == nil {

		conn.openStreamRW.Lock()
		conn.openStreamCount++
		conn.openStreamRW.Unlock()

		return &stdQuicStream{Stream: stream, conn: conn}, nil
	}

	conn.openStreamRW.Lock()
	conn.openStreamCount = 0
	conn.openStreamRW.Unlock()
	conn.close()
	return nil, err
}

func (conn *stdQuicConn) AcceptStream(ctx context.Context) (PeerStream, error) {

	stream, err := conn.Connection.AcceptStream(ctx)
	if err == nil {
		return stream, err
	}

	conn.close()
	return stream, err
}

func (conn *stdQuicConn) CloseWithError(errCode quic.ApplicationErrorCode, msg string) error {
	defer conn.close()
	return conn.Connection.CloseWithError(errCode, msg)
}

func (conn *stdQuicConn) Close() error {
	return conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (conn *stdQuicConn) close() {
	conn.rw.Lock()
	defer conn.rw.Unlock()
	if conn.closed {
		return
	}
	conn.closed = true
}

func (conn *stdQuicConn) isClosed() bool {
	conn.rw.RLock()
	defer conn.rw.RUnlock()
	return conn.closed
}

func (conn *stdQuicConn) detect() bool {
	conn.detectLocker.Lock()
	defer conn.detectLocker.Unlock()

	closed := conn.isClosed()
	if closed {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), QUIC_CONN_DETECT_TIMEOUT)
	defer cancel()

	_, err := conn.ReceiveDatagram(ctx)
	if err == nil || errors.Is(err, ctx.Err()) {
		return true
	}

	return false
}

func parseQuicConn(conn PeerConn) (*stdQuicConn, error) {
	quicConn, ok := conn.(*stdQuicConn)
	if !ok {
		return nil, ErrQuicConnInvalidPeerConn
	}
	return quicConn, nil
}

func extractPeerIDFromQuicConnection(conn quic.Connection) (PeerID, error) {
	connectionState := conn.ConnectionState()
	certificate := connectionState.TLS.PeerCertificates[0]
	return x509.MarshalPKIXPublicKey(certificate.PublicKey)
}

func releaseOpenStreamForQuicConn(conn *stdQuicConn) {
	conn.openStreamRW.Lock()
	defer conn.openStreamRW.Unlock()
	if conn.openStreamCount > 0 {
		conn.openStreamCount--
	}
}

func isIdleQuicConn(conn *stdQuicConn) bool {
	conn.openStreamRW.RLock()
	defer conn.openStreamRW.RUnlock()
	return conn.openStreamCount == 0
}
