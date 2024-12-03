package quic

import (
	"errors"
	"pan/app/peer"
	"slices"
	"sync"
)

var ErrQuicPeerRouteDuplicateAddress = errors.New("quic.PeerRoute Error: duplicate address")

type quicRoute struct {
	peerId peer.PeerID
	addrs  []string
	rw     sync.RWMutex
	sync.Mutex
}

func (qr *quicRoute) PeerID() peer.PeerID {
	return qr.peerId
}

func (qr *quicRoute) Addrs() []string {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return slices.Clone(qr.addrs)
}

func (qr *quicRoute) Available() bool {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return len(qr.addrs) > 0
}

func (qr *quicRoute) Contains(addr string) bool {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return slices.Contains(qr.addrs, addr)
}

func (qr *quicRoute) Store(addr string) error {
	qr.rw.Lock()
	defer qr.rw.Unlock()
	if ok := slices.Contains(qr.addrs, addr); ok {
		return ErrQuicPeerRouteDuplicateAddress
	}
	qr.addrs = append(qr.addrs, addr)
	return nil
}

func (qr *quicRoute) Delete(addr string) {
	qr.rw.Lock()
	defer qr.rw.Unlock()
	idx := slices.Index(qr.addrs, addr)
	if idx < 0 {
		return
	}
	qr.addrs = slices.Delete(qr.addrs, idx, idx+1)
}

// func (qr *quicPeerRoute) Dial(ctx context.Context) (quic.Connection, error) {

// 	qr.failureLocker.RLock()
// 	if qr.failures >= 3 {
// 		qr.Close()
// 		return nil, ErrQuicPeerRouteInvalid
// 	}
// 	qr.failureLocker.RUnlock()

// 	conn, err := qr.quicPeerModule.Dial(ctx, qr.address)
// 	if err == nil {
// 		peerId, err := parsePeerID(conn)
// 		if err == nil && !bytes.Equal(peerId, qr.peerId) {
// 			defer qr.Close()
// 			err = ErrQuicPeerRouteConflict
// 		}
// 		if err != nil {
// 			conn = nil
// 			defer conn.CloseWithError(quic.ApplicationErrorCode(quic.NoError), "")
// 		}
// 	}

// 	qr.failureLocker.Lock()
// 	defer qr.failureLocker.Unlock()
// 	if err != nil {
// 		qr.failures++
// 		failures := qr.failures
// 		if failures >= 3 {
// 			qr.Close()
// 		}
// 	} else {
// 		qr.failures = 0
// 	}

// 	return conn, err
// }
