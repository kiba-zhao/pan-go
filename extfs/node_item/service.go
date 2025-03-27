// Define node item service
package nodeitem

import (
	"errors"
	"os"
)

const (
	FileTypeFolder = "D"
	FileTypeFile   = "F"
)

type NodeItemInternalService interface {
	// TraverseAll iterates over all NodeItems in the repository, applying the given
	TraverseAll(func(NodeItem) error) error
	// Select a node item by its id
	Select(uint) (NodeItem, error)
	// SelectByName retrieves a node item by its name
	SelectByName(string) (NodeItem, error)
	// IsNotExist returns true if the given error is a ErrNodeItemNotFound error
	IsNotExist(error) bool
	// SelectAllWithEnabled retrieves all NodeItems from the repository which have the Enabled field matching the given argument.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	SelectAllWithEnabled(bool) ([]NodeItem, error)
}

type NodeItemService struct {
	NodeItemRepo NodeItemRepository
}

// IsNotExist returns true if the given error is a ErrNodeItemNotFound error
func (s *NodeItemService) IsNotExist(err error) bool {
	return errors.Is(err, ErrNodeItemNotFound)
}

// SelectAll returns a list of all NodeItems in the repository.
// The total count of NodeItems is returned as the first argument.
// If an error occurs during iteration, the iteration stops and
// the error is returned.
func (s *NodeItemService) SelectAll() (int64, []NodeItem, error) {
	items := make([]NodeItem, 0)
	err := s.TraverseAll(func(nodeItem NodeItem) error {
		items = append(items, nodeItem)
		return nil
	})

	if err != nil {
		return 0, nil, err
	}
	return int64(len(items)), items, nil
}

// Create creates a new NodeItem based on the provided fields.
// It checks the file system for the existence of the specified file path
// and determines the availability and type (file or folder) of the NodeItem.
// The NodeItem is then saved to the repository.
// Returns the created NodeItem and an error if any occurred during the process.

func (s *NodeItemService) Create(fields NodeItemFields) (NodeItem, error) {
	stat, err := os.Stat(fields.FilePath)
	if err != nil {
		return NodeItem{}, err
	}

	var nodeItem NodeItem
	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	nodeItem.Available = true
	if stat.IsDir() {
		nodeItem.FileType = FileTypeFolder
	} else {
		nodeItem.FileType = FileTypeFile
	}

	nodeItem_, err := s.NodeItemRepo.Save(nodeItem)
	if err == nil {
		setNodeItemAvailableWithFileStat(&nodeItem_)
	}
	return nodeItem_, err
}

// Update updates the NodeItem with the provided id based on the given fields.
// It checks the file system for the existence of the specified file path
// and determines the availability and type (file or folder) of the NodeItem.
// The NodeItem is then saved to the repository.
// Returns the updated NodeItem and an error if any occurred during the process.
func (s *NodeItemService) Update(fields NodeItemFields, id uint) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}

	stat, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return nodeItem, err
	}
	err = nil

	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	if stat.IsDir() {
		nodeItem.FileType = FileTypeFolder
	} else {
		nodeItem.FileType = FileTypeFile
	}

	nodeItem, err = s.NodeItemRepo.Save(nodeItem)
	if err == nil {
		setNodeItemAvailableWithFileStat(&nodeItem)
	}
	return nodeItem, err
}

// Select a node item by its id.
// If the item is not found, ErrNodeItemNotFound will be returned.
// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
// If the file system is unavailable, ErrNodeItemFileNotAvailable will be returned.
func (s *NodeItemService) Select(id uint) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	return nodeItem, nil
}

// SelectByName retrieves a NodeItem by its name from the repository.
// If the NodeItem is found, its availability and file statistics are updated.
// Returns the NodeItem and an error if any issues occur during retrieval.
func (s *NodeItemService) SelectByName(name string) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.SelectByName(name)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	return nodeItem, nil
}

// TraverseAll iterates over all NodeItems in the repository, applying the given
// function to each item. The availability and file statistics of each item are
// updated before being passed to the given function.
// If an error occurs during iteration, the iteration stops and
// the error is returned.
func (s *NodeItemService) TraverseAll(traverseFn func(NodeItem) error) error {
	return s.NodeItemRepo.TraverseAll(func(nodeItem NodeItem) error {
		setNodeItemAvailableWithFileStat(&nodeItem)
		return traverseFn(nodeItem)
	})
}

// SelectAllWithEnabled retrieves all NodeItems from the repository which have the Enabled field matching the given argument.
// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
func (s *NodeItemService) SelectAllWithEnabled(enabled bool) ([]NodeItem, error) {
	return s.NodeItemRepo.SelectAllWithEnabled(enabled)
}

// Delete removes a NodeItem from the repository by its id.
// If the NodeItem does not exist, ErrNodeItemNotFound will be returned.
// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
// Returns an error if any issues occur during deletion.
func (s *NodeItemService) Delete(id uint) error {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return err
	}
	err = s.NodeItemRepo.Delete(nodeItem)
	return err
}

func setNodeItemAvailableWithFileStat(nodeItem *NodeItem) {
	nodeItem.Available = *nodeItem.Enabled
	if !nodeItem.Available {
		return
	}

	stat, err := os.Stat(nodeItem.FilePath)
	nodeItem.Available = err == nil
	if !nodeItem.Available {
		return
	}

	nodeItem.Size = stat.Size()
	if stat.IsDir() {
		nodeItem.Available = nodeItem.FileType == FileTypeFolder
	} else {
		nodeItem.Available = nodeItem.FileType == FileTypeFile
	}
}
