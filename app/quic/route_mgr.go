package quic

import (
	"bytes"
	"slices"
	"sync"
)

type quicPeerRouteMgr struct {
	routes []*quicPeerRoute
	rw     sync.RWMutex
}

func (mgr *quicPeerRouteMgr) compareWithQuicRoute(route *quicPeerRoute, routeId []byte) int {
	return bytes.Compare(route.routeId, routeId)
}

func (mgr *quicPeerRouteMgr) Search(routeId []byte) *quicPeerRoute {
	mgr.rw.RLock()
	defer mgr.rw.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, routeId, mgr.compareWithQuicRoute)
	if !ok {
		return nil
	}
	return mgr.routes[idx]
}

func (mgr *quicPeerRouteMgr) SearchOrStore(route *quicPeerRoute) (*quicPeerRoute, bool) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.routeId, mgr.compareWithQuicRoute)
	if !ok {
		mgr.routes = slices.Insert(mgr.routes, idx, route)
		return route, false
	}
	return mgr.routes[idx], true
}

func (mgr *quicPeerRouteMgr) Delete(route *quicPeerRoute) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.routeId, mgr.compareWithQuicRoute)
	if !ok {
		return
	}
	mgr.routes = slices.Delete(mgr.routes, idx, idx+1)
}
