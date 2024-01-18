package peer

import (
	"cmp"
	"pan/memory"
)

type peerNodeItem struct {
	*memory.BucketItem[uint32]
	node Node
}

type peerNodeBucket = *memory.NestBucket[PeerId, uint32, *peerNodeItem]

type NodeManager struct {
	bucket *memory.Bucket[PeerId, peerNodeBucket]
}

// NewNodeManager ...
func NewNodeManager() *NodeManager {
	mgr := new(NodeManager)
	mgr.bucket = memory.NewBucket[PeerId, peerNodeBucket](comparePeerId)
	return mgr
}

// Count ...
func (mgr *NodeManager) Count(id PeerId) int {
	nodeBucket := mgr.bucket.GetItem(id)
	if nodeBucket == nil {
		return -1
	}
	return nodeBucket.Count()
}

// Get ...
func (mgr *NodeManager) Get(id PeerId) Node {
	nodeBucket := mgr.bucket.GetItem(id)
	if nodeBucket == nil {
		return nil
	}
	nodes := nodeBucket.GetAll()
	if nodes == nil || len(nodes) <= 0 {
		return nil
	}
	return nodes[0].node
}

// Save ...
func (mgr *NodeManager) Save(id PeerId, node Node) (*peerNodeItem, error) {
	nodeBucket := memory.NewNestBucket[PeerId, uint32, *peerNodeItem](id, cmp.Compare[uint32])
	nodeBucket, ok := mgr.bucket.GetOrAddItem(nodeBucket)

	nodeItem := new(peerNodeItem)
	nodeItem.node = node
	code := uint32(0)
	if ok {
		lastItem := nodeBucket.GetLastItem()
		if lastItem != nil {
			code = lastItem.HashCode() + 1
		}

	}
	nodeItem.BucketItem = memory.NewBucketItem[uint32](code)
	err := nodeBucket.AddItem(nodeItem)

	return nodeItem, err

}

// Delete ...
func (mgr *NodeManager) Delete(id PeerId, item *peerNodeItem) {
	nodeBucket := mgr.bucket.GetItem(id)
	if nodeBucket != nil {
		nodeBucket.RemoveItem(item)
	}
}
