package quic

import (
	"context"
	"pan/app/peer"
	"sync"

	"github.com/quic-go/quic-go"
)

type QuicConn interface {
	quic.Connection
	PeerID() peer.PeerID
	Available() bool
	Closed() bool
	CloseStream(quic.Stream)
}

type quicConn struct {
	quic.Connection
	peerId       peer.PeerID
	closed       bool
	closedLocker sync.Mutex
	mgr          *quicConnMgr
}

func (c *quicConn) PeerID() peer.PeerID {
	return c.peerId
}

func (c *quicConn) Available() bool {

	c.closedLocker.Lock()
	available := !c.closed
	c.closedLocker.Unlock()

	return available
}

func (c *quicConn) Closed() bool {

	c.closedLocker.Lock()
	defer c.closedLocker.Unlock()
	return c.closed
}

func (c *quicConn) CloseStream(stream quic.Stream) {
	stream.CancelRead(quic.StreamErrorCode(quic.NoError))
}

func (c *quicConn) AcceptStream(ctx context.Context) (quic.Stream, error) {
	stream, err := c.Connection.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}

	return &quicStream{Stream: stream, conn: c}, err
}

func (c *quicConn) OpenStream() (quic.Stream, error) {
	stream, err := c.Connection.OpenStream()
	if err != nil {
		return nil, err
	}

	return &quicStream{Stream: stream, conn: c, hangup: true}, err
}

func (c *quicConn) CloseWithError(code quic.ApplicationErrorCode, reason string) error {
	c.closedLocker.Lock()
	defer c.closedLocker.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true

	if c.mgr != nil {
		c.mgr.Delete(c)
	}
	return c.Connection.CloseWithError(code, reason)
}
