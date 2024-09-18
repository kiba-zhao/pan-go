package repositories

import "pan/extfs/models"

type NodeItemRepository interface {
	Save(models.NodeItem) (models.NodeItem, error)
	Select(uint) (models.NodeItem, error)
	Delete(models.NodeItem) error
	TraverseAll(func(models.NodeItem) error) error
}
