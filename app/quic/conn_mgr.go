// Define connection manager for quic
package quic

import (
	"bytes"
	"pan/app/peer"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

type quicConnMgr struct {
	connMatrix [][]QuicConn
	rw         sync.RWMutex
}

func (mgr *quicConnMgr) compare(connArr []QuicConn, peerId peer.PeerID) int {
	return bytes.Compare(connArr[0].PeerID(), peerId)
}

// Search returns a slice of QuicConn associated with the given peer ID if it exists.
// It acquires a read lock to ensure thread-safe access to the connections.
// If the peer ID does not exist, it returns nil.
func (mgr *quicConnMgr) Search(peerId peer.PeerID) []QuicConn {
	mgr.rw.RLock()
	defer mgr.rw.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, peerId, mgr.compare)
	if !ok {
		return nil
	}
	return slices.Clone(mgr.connMatrix[idx])
}

// SelectOrStore attempts to find the given QuicConn in the connection manager.
// If the connection does not exist, it inserts it into the connection matrix.
// It acquires a write lock to ensure thread-safe access to the connections.
// Returns the existing or newly inserted connection, and a boolean indicating
// if the connection was already present (true) or newly added (false).

func (mgr *quicConnMgr) SelectOrStore(conn QuicConn) (QuicConn, bool) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, conn.PeerID(), mgr.compare)
	if !ok {
		mgr.connMatrix = slices.Insert(mgr.connMatrix, idx, []QuicConn{conn})
		return conn, false
	}
	cidx := slices.Index(mgr.connMatrix[idx], conn)
	if cidx < 0 {
		mgr.connMatrix[idx] = append(mgr.connMatrix[idx], conn)
		return conn, false
	}
	return mgr.connMatrix[idx][cidx], true
}

// Delete removes the given QuicConn from the connection manager.
// It acquires a write lock to ensure thread-safe access to the connections.
// If the connection is not found, it returns immediately.
func (mgr *quicConnMgr) Delete(conn QuicConn) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, conn.PeerID(), mgr.compare)
	if !ok {
		return
	}

	cidx := slices.Index(mgr.connMatrix[idx], conn)
	if cidx < 0 {
		return
	}
	if len(mgr.connMatrix[idx]) <= 1 {
		mgr.connMatrix = slices.Delete(mgr.connMatrix, idx, idx+1)
		return
	}
	mgr.connMatrix[idx] = slices.Delete(mgr.connMatrix[idx], cidx, cidx+1)

}

// Clean closes all connections associated with the given peer ID.
// It acquires a write lock to ensure thread-safe access to the connections.
func (mgr *quicConnMgr) Clean(peer peer.PeerID) {
	var connArr []QuicConn
	mgr.rw.Lock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, peer, mgr.compare)
	if ok {
		connArr = mgr.connMatrix[idx]
		mgr.connMatrix = slices.Delete(mgr.connMatrix, idx, idx+1)
	}
	mgr.rw.Unlock()

	if len(connArr) <= 0 {
		return
	}
	for _, conn := range connArr {
		conn.CloseWithError(quic.ApplicationErrorCode(0), "")
	}
}
