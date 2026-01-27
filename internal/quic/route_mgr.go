// Define route manager for quic
package quic

import (
	"bytes"
	"pan/internal/peer"
	"slices"
	"sync"
)

type stdQuicRouteMgr struct {
	routes []*stdQuicRoute
	rw     sync.RWMutex
}

func (mgr *stdQuicRouteMgr) compare(route *stdQuicRoute, peerId peer.PeerID) int {
	return bytes.Compare(route.peerId, peerId)
}

// Search returns the route associated with the given peer ID if it exists.
// It acquires a read lock to ensure thread-safe access to the routes.
// If the route does not exist, it returns nil.
func (mgr *stdQuicRouteMgr) Search(peerId peer.PeerID) *stdQuicRoute {
	mgr.rw.RLock()
	defer mgr.rw.RUnlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, peerId, mgr.compare)
	if !ok {
		return nil
	}
	return mgr.routes[idx]
}

// SearchOrStore searches for a route in the route manager, and if it does not exist, stores the given route.
// It acquires a write lock to ensure thread-safe access to the routes.
// If the route does not exist, it returns the given route with the second return value set to false.
// If the route exists, it returns the existing route with the second return value set to true.
func (mgr *stdQuicRouteMgr) SearchOrStore(route *stdQuicRoute) (*stdQuicRoute, bool) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.PeerID(), mgr.compare)
	if !ok {
		mgr.routes = slices.Insert(mgr.routes, idx, route)
		return route, false
	}
	return mgr.routes[idx], true
}

// Delete removes the given route from the stdQuicRouteMgr.
//
// It acquires a write lock to ensure thread-safe access to the routes.
// If the route is not found, it returns immediately.
func (mgr *stdQuicRouteMgr) Delete(route *stdQuicRoute) {
	mgr.rw.Lock()
	defer mgr.rw.Unlock()
	idx, ok := slices.BinarySearchFunc(mgr.routes, route.PeerID(), mgr.compare)
	if !ok {
		return
	}
	mgr.routes = slices.Delete(mgr.routes, idx, idx+1)
}
