// Define route for quic
package quic

import (
	"errors"
	"pan/lib/peer"
	"slices"
	"sync"
)

var ErrQuicPeerRouteDuplicateAddress = errors.New("quic.PeerRoute Error: duplicate address")

type stdQuicRoute struct {
	peerId peer.PeerID
	addrs  []string
	rw     sync.RWMutex
	sync.Mutex
}

// PeerID returns the peer ID of the peer route.
func (qr *stdQuicRoute) PeerID() peer.PeerID {
	return qr.peerId
}

// Addrs returns a copy of the list of addresses associated with the stdQuicRoute.
// It acquires a read lock to ensure thread-safe access to the addresses.

func (qr *stdQuicRoute) Addrs() []string {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return slices.Clone(qr.addrs)
}

// Available returns true if the stdQuicRoute has at least one associated address.
// It acquires a read lock to ensure thread-safe access to the addresses.
func (qr *stdQuicRoute) Available() bool {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return len(qr.addrs) > 0
}

// Contains returns true if the given address is associated with the stdQuicRoute.
// It acquires a read lock to ensure thread-safe access to the addresses.
func (qr *stdQuicRoute) Contains(addr string) bool {
	qr.rw.RLock()
	defer qr.rw.RUnlock()
	return slices.Contains(qr.addrs, addr)
}

// Store adds the given address to the stdQuicRoute.
//
// It acquires a write lock to ensure thread-safe access to the addresses.
// If the address is already associated with the stdQuicRoute, it returns ErrQuicPeerRouteDuplicateAddress.
func (qr *stdQuicRoute) Store(addr string) error {
	qr.rw.Lock()
	defer qr.rw.Unlock()
	if ok := slices.Contains(qr.addrs, addr); ok {
		return ErrQuicPeerRouteDuplicateAddress
	}
	qr.addrs = append(qr.addrs, addr)
	return nil
}

// Delete removes the given address from the stdQuicRoute.
//
// It acquires a write lock to ensure thread-safe access to the addresses.
// If the address is not found, the function returns immediately.

func (qr *stdQuicRoute) Delete(addr string) {
	qr.rw.Lock()
	defer qr.rw.Unlock()
	idx := slices.Index(qr.addrs, addr)
	if idx < 0 {
		return
	}
	qr.addrs = slices.Delete(qr.addrs, idx, idx+1)
}
