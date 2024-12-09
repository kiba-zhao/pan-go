package quic

import (
	"cmp"
	"context"
	"pan/app/peer"
	"pan/logger"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

const QuicConnAvaliableWindowSize = (1 << 10) * 512

type QuicConn interface {
	quic.Connection
	PeerID() peer.PeerID
	Available() bool
	Closed() bool
	CloseStream(quic.Stream)
	OnStreamRead(quic.StreamID, int)
	OnStreamWrite(quic.StreamID, int)
}

type quicStreamWindow struct {
	streamId          quic.StreamID
	readBytes         int
	writeBytes        int
	readRW            sync.RWMutex
	writeRW           sync.RWMutex
	blockedReadBytes  int
	blockedWriteBytes int
	blockedReadRW     sync.RWMutex
	blockedWriteRW    sync.RWMutex
}

type quicConn struct {
	quic.Connection
	peerId          peer.PeerID
	closed          bool
	closedLocker    sync.Mutex
	mgr             *quicConnMgr
	streamWindows   []*quicStreamWindow
	streamWindowsRW sync.RWMutex
	agent           *quicPeerAgent
	syncCount       uint8
	syncLocker      sync.Mutex
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
	if len(c.streamWindows) > 1 {
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

	if size > 0 {
		window.blockedReadRW.Lock()
		window.blockedReadBytes += size
		window.blockedReadRW.Unlock()

		syncStreamBytes(c)
	}
}

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

	if size > 0 {
		window.blockedWriteRW.Lock()
		window.blockedWriteBytes += size
		window.blockedWriteRW.Unlock()

		syncStreamBytes(c)
	}
}

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

func syncStreamBytes(c *quicConn) {

	c.syncLocker.Lock()

	if c.syncCount > 1 {
		c.syncLocker.Unlock()
		return
	}
	syncCount := c.syncCount
	c.syncCount++
	c.syncLocker.Unlock()

	if syncCount < 1 {
		go syncWorker(c)
	}

}

func syncWorker(c *quicConn) {
	defer completeSyncWorker(c)

	if c.Closed() {
		return
	}

	c.streamWindowsRW.RLock()
	if len(c.streamWindows) <= 0 {
		c.streamWindowsRW.RUnlock()
		return
	}
	windows := slices.Clone(c.streamWindows)
	c.streamWindowsRW.RUnlock()

	var list QuicStreamBytesList
	for _, window := range windows {
		var stream QuicStreamBytes
		stream.StreamID = int64(window.streamId)

		window.blockedReadRW.Lock()
		stream.ReadSize = int32(window.blockedReadBytes)
		window.blockedReadBytes = 0
		window.blockedReadRW.Unlock()

		window.blockedWriteRW.Lock()
		stream.WriteSize = int32(window.blockedWriteBytes)
		window.blockedWriteBytes = 0
		window.blockedWriteRW.Unlock()
		if stream.ReadSize == 0 && stream.WriteSize == 0 {
			continue
		}
		list.Streams = append(list.Streams, &stream)
	}

	if len(list.Streams) <= 0 {
		return
	}

	err := c.agent.Sync(c, &list)
	if err != nil {
		logger.Default().Log(context.Background(), logger.LevelError, err.Error())
	}

}

func completeSyncWorker(c *quicConn) {
	c.syncLocker.Lock()
	defer c.syncLocker.Unlock()

	c.syncCount--
	if c.syncCount > 0 {
		go syncWorker(c)
	}
}
