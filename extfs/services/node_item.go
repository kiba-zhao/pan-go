package services

import (
	"os"
	"pan/app/constant"
	"pan/extfs/models"
	"pan/extfs/repositories"
)

type NodeItemService struct {
	NodeItemRepo repositories.NodeItemRepository
}

func (s *NodeItemService) Create(fields models.NodeItemFields) (models.NodeItem, error) {
	_, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return models.NodeItem{}, err
	}
	available := err == nil

	var nodeItem models.NodeItem
	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled
	nodeItem.Available = available

	nodeItem_, err := s.NodeItemRepo.Save(nodeItem)
	if err == nil {
		setNodeItemAvailable(&nodeItem_)
	}
	return nodeItem_, err
}

func (s *NodeItemService) Update(fields models.NodeItemFields, id uint) (models.NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}

	_, err = os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return nodeItem, err
	}
	err = nil

	nodeItem.Name = fields.Name
	nodeItem.FilePath = fields.FilePath
	nodeItem.Enabled = fields.Enabled

	nodeItem, err = s.NodeItemRepo.Save(nodeItem)
	if err == constant.ErrNotFound {
		return nodeItem, constant.ErrConflict
	}
	if err == nil {
		setNodeItemAvailable(&nodeItem)
	}
	return nodeItem, err
}

func (s *NodeItemService) Select(id uint) (models.NodeItem, error) {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return nodeItem, err
	}
	setNodeItemAvailable(&nodeItem)
	return nodeItem, nil
}

func (s *NodeItemService) Delete(id uint) error {
	nodeItem, err := s.NodeItemRepo.Select(id)
	if err != nil {
		return err
	}
	err = s.NodeItemRepo.Delete(nodeItem)
	if err == constant.ErrNotFound {
		return constant.ErrConflict
	}
	return err
}

func setNodeItemAvailable(nodeItem *models.NodeItem) {
	nodeItem.Available = *nodeItem.Enabled
	if nodeItem.Available {
		_, err := os.Stat(nodeItem.FilePath)
		nodeItem.Available = err == nil
	}
}
