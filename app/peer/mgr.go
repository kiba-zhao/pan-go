package peer

import (
	"bytes"
	"slices"
	"sync"
)

type PeerManager interface {
	TraversePeerID(func(PeerID) error) error
	TraversePeerNode(PeerID, func(PeerNode) bool)
	Search(PeerID) []PeerNode
	Delete(PeerNode)
	SearchOrStore(PeerNode) (PeerNode, bool)
	Count(PeerID) int
}

type peerManager struct {
	locker sync.RWMutex
	matrix [][]PeerNode
}

func (mgr *peerManager) TraversePeerID(traverse func(PeerID) error) error {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	for _, peerNodes := range mgr.matrix {
		if len(peerNodes) > 0 {
			if err := traverse(peerNodes[0].PeerID()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mgr *peerManager) compareWithPeerID(peerNodes []PeerNode, peerId PeerID) int {
	return bytes.Compare(peerNodes[0].PeerID(), peerId)
}

func (mgr *peerManager) compareWithResourceID(peerNode PeerNode, resourceId []byte) int {
	return bytes.Compare(peerNode.PeerResourceID(), resourceId)
}

func (mgr *peerManager) TraversePeerNode(peerId PeerID, traverse func(PeerNode) bool) {

	peerNodes := mgr.Search(peerId)
	if len(peerNodes) <= 0 {
		return
	}

	for _, peerNode := range peerNodes {
		if !traverse(peerNode) {
			break
		}
	}

}

func (mgr *peerManager) Search(peerId PeerID) []PeerNode {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.matrix, peerId, mgr.compareWithPeerID)
	if ok {
		return slices.Clone(mgr.matrix[idx])
	}
	return nil
}

func (mgr *peerManager) Delete(peerNode PeerNode) {
	mgr.locker.Lock()
	defer mgr.locker.Unlock()

	idx, ok := slices.BinarySearchFunc(mgr.matrix, peerNode.PeerID(), mgr.compareWithPeerID)
	if !ok {
		return
	}

	idx_, ok := slices.BinarySearchFunc(mgr.matrix[idx], peerNode.PeerResourceID(), mgr.compareWithResourceID)
	if !ok {
		return
	}

	if len(mgr.matrix[idx]) <= 1 {
		mgr.matrix = slices.Delete(mgr.matrix, idx, idx+1)
		return
	}
	mgr.matrix[idx] = slices.Delete(mgr.matrix[idx], idx_, idx_+1)
}

func (mgr *peerManager) SearchOrStore(peerNode PeerNode) (PeerNode, bool) {

	mgr.locker.Lock()
	defer mgr.locker.Unlock()

	idx, ok := slices.BinarySearchFunc(mgr.matrix, peerNode.PeerID(), mgr.compareWithPeerID)

	if !ok {
		mgr.matrix = slices.Insert(mgr.matrix, idx, []PeerNode{peerNode})
		return peerNode, false
	}

	idx_, ok := slices.BinarySearchFunc(mgr.matrix[idx], peerNode.PeerResourceID(), mgr.compareWithResourceID)
	if !ok {
		mgr.matrix[idx] = slices.Insert(mgr.matrix[idx], idx_, peerNode)
		return peerNode, false
	}

	return mgr.matrix[idx][idx_], true

}

func (mgr *peerManager) Count(peerId PeerID) int {

	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.matrix, peerId, mgr.compareWithPeerID)
	if ok {
		return len(mgr.matrix[idx])
	}
	return 0
}
