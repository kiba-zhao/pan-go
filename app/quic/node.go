package quic

import (
	"context"
	"io"
	"pan/app/peer"
	"sync"

	"github.com/quic-go/quic-go"
)

type quicPeerNode struct {
	resourceId     peer.PeerResourceID
	nodeId         peer.PeerID
	conn           quic.Connection
	quicPeerModule QuicPeerModule
	mgr            peer.PeerManager
	streamCount    int8
	rw             sync.RWMutex
}

func (qn *quicPeerNode) PeerID() peer.PeerID {
	return qn.nodeId
}

func (qn *quicPeerNode) PeerType() peer.PeerType {
	return peer.PeerTypeAlive
}

func (qn *quicPeerNode) Do(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
	qn.increaseStream()
	stream, err := qn.quicPeerModule.Do(ctx, qn.conn, reader)
	if err != nil {
		qn.decreaseStream()
		return nil, err
	}
	return &quicPeerStream{Stream: stream, quicPeerNode: qn}, err
}

func (qn *quicPeerNode) Close() error {
	qn.rw.Lock()
	qn.streamCount = -1
	qn.rw.Unlock()

	qn.mgr.Delete(qn)
	return qn.conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
}

func (qn *quicPeerNode) PeerResourceID() peer.PeerResourceID {
	return qn.resourceId
}

func (qn *quicPeerNode) IsIdle() bool {
	qn.rw.RLock()
	defer qn.rw.RUnlock()
	return qn.streamCount >= 0 && qn.streamCount < 2
}

func (qn *quicPeerNode) increaseStream() {
	qn.rw.Lock()
	defer qn.rw.Unlock()
	qn.streamCount++
}

func (qn *quicPeerNode) decreaseStream() {
	qn.rw.Lock()
	defer qn.rw.Unlock()
	if qn.streamCount <= 0 {
		return
	}
	qn.streamCount--
}
