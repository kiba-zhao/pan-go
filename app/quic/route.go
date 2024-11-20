package quic

import (
	"bytes"
	"context"
	"errors"
	"io"
	"pan/app/peer"
	"sync"

	"github.com/quic-go/quic-go"
)

var ErrQuicPeerRouteInvalid = errors.New("quic.PeerRoute Error: Invalid Route")
var ErrQuicPeerRouteConflict = errors.New("quic.PeerRoute Error:Peer Conflict")

type quicPeerRoute struct {
	resourceId     peer.PeerResourceID
	peerId         peer.PeerID
	address        string
	quicPeerModule QuicPeerModule
	failures       uint8
	failureLocker  sync.RWMutex
	routeId        []byte
	closed         bool
	closedRW       sync.RWMutex
	mgr            *quicPeerRouteMgr
}

func (qr *quicPeerRoute) PeerID() peer.PeerID {
	return qr.peerId
}

func (qr *quicPeerRoute) PeerType() peer.PeerType {
	return peer.PeerTypeReachable
}

func (qr *quicPeerRoute) Dial(ctx context.Context) (quic.Connection, error) {

	qr.failureLocker.RLock()
	if qr.failures >= 3 {
		qr.Close()
		return nil, ErrQuicPeerRouteInvalid
	}
	qr.failureLocker.RUnlock()

	conn, err := qr.quicPeerModule.Dial(ctx, qr.address)
	if err == nil {
		peerId, err := parsePeerID(conn)
		if err == nil && !bytes.Equal(peerId, qr.peerId) {
			defer qr.Close()
			err = ErrQuicPeerRouteConflict
		}
		if err != nil {
			conn = nil
			defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
		}
	}

	qr.failureLocker.Lock()
	defer qr.failureLocker.Unlock()
	if err != nil {
		qr.failures++
		failures := qr.failures
		if failures >= 3 {
			qr.Close()
		}
	} else {
		qr.failures = 0
	}

	return conn, err
}

func (qr *quicPeerRoute) Do(ctx context.Context, reader io.Reader) (io.ReadCloser, error) {
	conn, err := qr.Dial(ctx)
	if err != nil {
		return nil, err
	}

	node, err := qr.quicPeerModule.Serve(conn)
	if err != nil {
		return nil, err
	}
	return node.Do(ctx, reader)
}

func (qr *quicPeerRoute) Close() error {
	qr.closedRW.Lock()
	qr.closed = true
	qr.closedRW.Unlock()

	qr.mgr.Delete(qr)
	return nil
}

func (qr *quicPeerRoute) PeerResourceID() peer.PeerResourceID {
	return qr.resourceId
}

func (qr *quicPeerRoute) IsIdle() bool {
	qr.closedRW.RLock()
	defer qr.closedRW.RUnlock()
	return !qr.closed
}
