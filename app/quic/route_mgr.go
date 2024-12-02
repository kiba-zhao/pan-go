package quic

import (
	"bytes"
	"pan/app/peer"
	"slices"
	"sync"
)

type quicRouteMgr struct {
	routes []*quicRoute
	rw     sync.RWMutex
}

func (mgr *quicRouteMgr) compare(route *quicRoute, peerId peer.PeerID) int {
	return bytes.Compare(route.peerId, peerId)
}

func (mgr *quicRouteMgr) Search(peerId peer.PeerID) *quicRoute {
	mgr.rw.RLock()
	defer mgr.rw.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, peerId, mgr.compare)
	if !ok {
		return nil
	}
	return mgr.routes[idx]
}

func (mgr *quicRouteMgr) SearchOrStore(route *quicRoute) (*quicRoute, bool) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.PeerID(), mgr.compare)
	if !ok {
		mgr.routes = slices.Insert(mgr.routes, idx, route)
		return route, false
	}
	return mgr.routes[idx], true
}

func (mgr *quicRouteMgr) Delete(route *quicRoute) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.PeerID(), mgr.compare)
	if !ok {
		return
	}
	mgr.routes = slices.Delete(mgr.routes, idx, idx+1)
}
