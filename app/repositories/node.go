package repositories

import "pan/app/models"

type NodeRepository interface {
	Search(models.NodeSearchCondition) (int64, []models.Node, error)
	Save(models.Node) (models.Node, error)
	Select(uint) (models.Node, error)
	Delete(models.Node) error
	SelectByNodeID(string) (models.Node, error)
	TraverseWithNodeIDs(func(models.Node) error, []string) error
}
