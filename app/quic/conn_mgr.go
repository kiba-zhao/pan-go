package quic

import (
	"bytes"
	"pan/app/peer"
	"slices"
	"sync"

	"github.com/quic-go/quic-go"
)

type QuicConnMgr interface {
	Search(peerId peer.PeerID) []QuicConn
	SelectOrStore(conn QuicConn) (QuicConn, bool)
	Delete(conn QuicConn)
	Clean(peer peer.PeerID)
}

type quicConnMgr struct {
	connMatrix [][]QuicConn
	rw         sync.RWMutex
}

func (mgr *quicConnMgr) compare(connArr []QuicConn, peerId peer.PeerID) int {
	return bytes.Compare(connArr[0].PeerID(), peerId)
}

func (mgr *quicConnMgr) Search(peerId peer.PeerID) []QuicConn {
	mgr.rw.RLock()
	defer mgr.rw.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, peerId, mgr.compare)
	if !ok {
		return nil
	}
	return slices.Clone(mgr.connMatrix[idx])
}

func (mgr *quicConnMgr) SelectOrStore(conn QuicConn) (QuicConn, bool) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.connMatrix, conn.PeerID(), mgr.compare)
	if !ok {
		conn.AssignManager(mgr)
		mgr.connMatrix = slices.Insert(mgr.connMatrix, idx, []QuicConn{conn})
		return conn, false
	}
	cidx := slices.Index(mgr.connMatrix[idx], conn)
	if cidx < 0 {
		conn.AssignManager(mgr)
		mgr.connMatrix[idx] = append(mgr.connMatrix[idx], conn)
		return conn, false
	}
	return mgr.connMatrix[idx][cidx], true
}

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
