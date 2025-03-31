// Define connection for quic
package quic

import (
	"cmp"
	"context"
	"pan/lib/peer"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

const QuicConnAvaliableWindowSize = (1 << 10) * 512

// Quic Connection Interface
type QuicConn interface {
	// Extends quic.Connection interface
	quic.Connection
	// PeerID returns the peer ID of the connection.
	PeerID() peer.PeerID
	// Available returns true if the connection is available.
	Available() bool
	// Closed returns true if the connection is closed.
	Closed() bool
	// CloseWithError closes the connection with the given error code and reason.
	CloseStream(quic.Stream)
}

type quicStreamWindow struct {
	streamId   quic.StreamID
	readBytes  int
	writeBytes int
	readRW     sync.RWMutex
	writeRW    sync.RWMutex
}

type quicConn struct {
	quic.Connection
	peerId          peer.PeerID
	closed          bool
	closedLocker    sync.Mutex
	mgr             *quicConnMgr
	streamWindows   []*quicStreamWindow
	streamWindowsRW sync.RWMutex
}

func (c *quicConn) PeerID() peer.PeerID {
	return c.peerId
}

func (c *quicConn) Available() bool {
	available := !c.Closed()
	if !available {
		return available
	}

	c.streamWindowsRW.RLock()
	if len(c.streamWindows) > 2 {
		readBytes := 0
		writeBytes := 0
		for _, window := range c.streamWindows {

			window.writeRW.RLock()
			writeBytes += window.writeBytes
			window.writeRW.RUnlock()
			if writeBytes > QuicConnAvaliableWindowSize || writeBytes < -1*QuicConnAvaliableWindowSize {
				available = false
				break
			}

			window.readRW.RLock()
			readBytes += window.readBytes
			window.readRW.RUnlock()
			if readBytes < -1*QuicConnAvaliableWindowSize || readBytes > QuicConnAvaliableWindowSize {
				available = false
				break
			}

		}
	}
	c.streamWindowsRW.RUnlock()
	return available
}

func (c *quicConn) Closed() bool {

	c.closedLocker.Lock()
	defer c.closedLocker.Unlock()
	return c.closed
}

func (c *quicConn) CloseStream(stream quic.Stream) {
	stream.CancelRead(quic.StreamErrorCode(quic.NoError))
	c.streamWindowsRW.Lock()
	c.streamWindows = removeStreamWindow(c.streamWindows, stream.StreamID())
	c.streamWindowsRW.Unlock()
}

// OnStreamRead updates the read byte count for a specific stream.
// It locates the stream's window via its streamId and increments
// the readBytes by the specified size, provided the connection
// is not closed and size is greater than zero. If the stream window
// is not found, the function exits without making changes.
// It implements the OnStreamRead method from the quic.Connection interface.

func (c *quicConn) OnStreamRead(streamId quic.StreamID, size int) {
	if c.Closed() {
		return
	}
	if size == 0 {
		return
	}
	c.streamWindowsRW.RLock()
	window := searchStreamWindow(c.streamWindows, streamId)
	c.streamWindowsRW.RUnlock()

	if window == nil {
		return
	}
	window.readRW.Lock()
	window.readBytes += size
	window.readRW.Unlock()

}

// OnStreamWrite updates the write byte count for a specific stream.
// It locates the stream's window via its streamId and increments
// the writeBytes by the specified size, provided the connection
// is not closed and size is greater than zero. If the stream window
// is not found, the function exits without making changes.
// It implements the OnStreamWrite method from the quic.Connection interface.
func (c *quicConn) OnStreamWrite(streamId quic.StreamID, size int) {
	if c.Closed() {
		return
	}
	if size == 0 {
		return
	}
	c.streamWindowsRW.RLock()
	window := searchStreamWindow(c.streamWindows, streamId)
	c.streamWindowsRW.RUnlock()
	if window == nil {
		return
	}
	window.writeRW.Lock()
	window.writeBytes += size
	window.writeRW.Unlock()
}

// AcceptStream implements the AcceptStream method of the quic.Connection interface.
// It adds the new stream to the connection's stream window list and returns a
// quicStream object wrapping the new stream and the connection.
func (c *quicConn) AcceptStream(ctx context.Context) (quic.Stream, error) {
	stream, err := c.Connection.AcceptStream(ctx)
	if err != nil {
		return nil, err
	}

	c.streamWindowsRW.Lock()
	c.streamWindows = storeStreamWindow(c.streamWindows, &quicStreamWindow{streamId: stream.StreamID()})
	c.streamWindowsRW.Unlock()
	return &quicStream{Stream: stream, conn: c}, err
}

// OpenStream implements the OpenStream method of the quic.Connection interface.
// It opens a new stream via the underlying quic.Connection and adds the new
// stream to the connection's stream window list. It then returns a quicStream
// object wrapping the new stream and the connection.
func (c *quicConn) OpenStream() (quic.Stream, error) {
	stream, err := c.Connection.OpenStream()
	if err != nil {
		return nil, err
	}

	c.streamWindowsRW.Lock()
	c.streamWindows = storeStreamWindow(c.streamWindows, &quicStreamWindow{streamId: stream.StreamID()})
	c.streamWindowsRW.Unlock()
	return &quicStream{Stream: stream, conn: c, hangup: true}, err
}

// CloseWithError closes the quicConn with the specified error code and reason.
// It first locks the closed state, checks if the connection is already closed,
// and if not, marks it as closed. It then removes the connection from its manager
// if applicable, and finally calls CloseWithError on the underlying connection.
// Returns an error if the underlying connection encounters an issue during closure.

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

func compareStreamWindow(window *quicStreamWindow, streamId quic.StreamID) int {
	return cmp.Compare(window.streamId, streamId)
}

func removeStreamWindow(streamWindows []*quicStreamWindow, streamId quic.StreamID) []*quicStreamWindow {
	idx, ok := slices.BinarySearchFunc(streamWindows, streamId, compareStreamWindow)
	if ok {
		streamWindows = slices.Delete(streamWindows, idx, idx+1)
	}
	return streamWindows
}

func searchStreamWindow(streamWindows []*quicStreamWindow, streamId quic.StreamID) *quicStreamWindow {
	idx, ok := slices.BinarySearchFunc(streamWindows, streamId, compareStreamWindow)
	if ok {
		return streamWindows[idx]
	}
	return nil
}

func storeStreamWindow(streamWindows []*quicStreamWindow, window *quicStreamWindow) []*quicStreamWindow {
	idx, ok := slices.BinarySearchFunc(streamWindows, window.streamId, compareStreamWindow)
	if ok {
		streamWindows[idx] = window
		return streamWindows
	}
	return slices.Insert(streamWindows, idx, window)
}
