package peer

import (
	"bytes"
	"pan/memory"
	"sync"
)

type peerRoute struct {
	*memory.BucketItem[[]byte]
	rw        *sync.RWMutex
	NodeType  NodeType
	Addr      []byte
	FailedNum uint8
}

type peerRouteBucket = *memory.NestBucket[PeerId, []byte, *peerRoute]

type RouteManager struct {
	bucket *memory.Bucket[PeerId, peerRouteBucket]
}

// NewRouteManager ...
func NewRouteManager() *RouteManager {
	mgr := new(RouteManager)
	mgr.bucket = memory.NewBucket[PeerId, peerRouteBucket](comparePeerId)
	return mgr
}

// Count ...
func (mgr *RouteManager) Count(id PeerId) int {
	routeBucket := mgr.bucket.GetItem(id)
	if routeBucket == nil {
		return -1
	}
	return routeBucket.Count()
}

// GetAll ...
func (mgr *RouteManager) GetAll(id PeerId) []*peerRoute {
	routeBucket := mgr.bucket.GetItem(id)
	if routeBucket == nil {
		return nil
	}
	routes := routeBucket.GetAll()
	if routes == nil || len(routes) <= 0 {
		return nil
	}
	return routes
}

// Save ...
func (mgr *RouteManager) Save(id PeerId, node Node) (*peerRoute, bool) {
	route := new(peerRoute)
	route.Addr = node.Addr()
	route.NodeType = node.Type()
	route.FailedNum = 0
	route.rw = new(sync.RWMutex)
	code := make([]byte, 0)
	code = append(code, route.NodeType)
	code = append(code, route.Addr...)
	route.BucketItem = memory.NewBucketItem[[]byte](code)

	routeBucket := memory.NewNestBucket[PeerId, []byte, *peerRoute](id, bytes.Compare)
	routeBucket, _ = mgr.bucket.GetOrAddItem(routeBucket)
	routeItem, ok := routeBucket.GetOrAddItem(route)
	if ok {
		routeItem.rw.Lock()
		routeItem.FailedNum = 0
		routeItem.rw.Unlock()
	}
	return routeItem, ok
}

// Delete ...
func (mgr *RouteManager) Delete(id PeerId, item *peerRoute) {
	routeBucket := mgr.bucket.GetItem(id)
	if routeBucket != nil {
		routeBucket.RemoveItem(item)
	}

}
