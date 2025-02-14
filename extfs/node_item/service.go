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
	TraverseAll(func(NodeItem) error) error
	Select(uint) (NodeItem, error)
	SelectByName(string) (NodeItem, error)
	IsNotExist(error) bool
	SelectAllWithEnabled(bool) ([]NodeItem, error)
}

type NodeItemService struct {
	NodeItemRepo NodeItemRepository
}

func (s *NodeItemService) IsNotExist(err error) bool {
	return errors.Is(err, ErrNodeItemNotFound)
}

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

func (s *NodeItemService) Create(fields NodeItemFields) (NodeItem, error) {
	stat, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return NodeItem{}, err
	}
	available := err == nil

	var nodeItem NodeItem
	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	nodeItem.Available = available
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

func (s *NodeItemService) Select(id uint) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	return nodeItem, nil
}

func (s *NodeItemService) SelectByName(name string) (NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.SelectByName(name)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailableWithFileStat(&nodeItem)
	return nodeItem, nil
}

func (s *NodeItemService) TraverseAll(traverseFn func(NodeItem) error) error {
	return s.NodeItemRepo.TraverseAll(func(nodeItem NodeItem) error {
		setNodeItemAvailableWithFileStat(&nodeItem)
		return traverseFn(nodeItem)
	})
}

func (s *NodeItemService) SelectAllWithEnabled(enabled bool) ([]NodeItem, error) {
	return s.NodeItemRepo.SelectAllWithEnabled(enabled)
}

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
