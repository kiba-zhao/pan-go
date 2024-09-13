package repositories

import "pan/extfs/models"

type NodeItemRepository interface {
	Save(nodeItem models.NodeItem) (models.NodeItem, error)
	Select(id uint) (models.NodeItem, error)
	Delete(nodeItem models.NodeItem) error
}
