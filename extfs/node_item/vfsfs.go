package nodeitem

type VFSFUSENodeItemRuntime interface {
	// NodeItemService returns the node item service.
	NodeItemService() *NodeItemService
}
