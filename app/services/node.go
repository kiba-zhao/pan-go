package services

import (
	"encoding/base64"
	"pan/app/models"
	"pan/app/node"
	"pan/app/repositories"
	"strings"
)

type NodeManagerProvider interface {
	NodeManager() node.NodeManager
}

type NodeService struct {
	NodeRepo repositories.NodeRepository
	Provider NodeManagerProvider
}

func (s *NodeService) Search(conditions models.NodeSearchCondition) (total int64, items []models.Node, err error) {
	total, items, err = s.NodeRepo.Search(conditions)
	if err != nil {
		return
	}

	if conditions.Blocked != nil && *conditions.Blocked {
		return
	}

	mgr := s.Provider.NodeManager()
	if mgr == nil {
		return
	}
	items_ := make([]models.Node, 0)
	for _, item := range items {
		if !item.Blocked {
			setNodeOnline(mgr, &item)
		}
		if conditions.Online != nil && *conditions.Online != item.Online {
			continue
		}
		items_ = append(items_, item)
	}
	items = items_
	return
}

func (s *NodeService) Select(id uint) (models.Node, error) {

	model, err := s.NodeRepo.Select(id)
	if err == nil && !model.Blocked {
		mgr := s.Provider.NodeManager()
		if mgr != nil {
			err = setNodeOnline(mgr, &model)
		}
	}
	return model, err
}

func (s *NodeService) Delete(id uint) error {
	model, err := s.NodeRepo.Select(id)
	if err != nil {
		return err
	}
	err = s.NodeRepo.Delete(model)
	if err != nil {
		return err
	}
	if !model.Blocked {
		mgr := s.Provider.NodeManager()
		if mgr != nil {
			err = closeNode(mgr, &model)
		}
	}
	return err
}

func (s *NodeService) Create(fields models.NodeFields) (models.Node, error) {

	var model models.Node
	model.Name = fields.Name
	model.NodeID = fields.NodeID
	model.Blocked = fields.Blocked != nil && *fields.Blocked

	return s.NodeRepo.Save(model)

}

func (s *NodeService) Update(id uint, fields models.NodeFields) (models.Node, error) {

	model, err := s.NodeRepo.Select(id)
	if err != nil {
		return model, err
	}

	dirty := false
	needClosed := false
	name := strings.Trim(fields.Name, " ")
	if len(name) > 0 {
		dirty = true
		model.Name = name
	}
	if fields.Blocked != nil {
		dirty = true
		model.Blocked = *fields.Blocked
		needClosed = *fields.Blocked
	}

	if dirty {
		model, err = s.NodeRepo.Save(model)
		if err == nil && needClosed {
			mgr := s.Provider.NodeManager()
			if mgr != nil {
				err = closeNode(mgr, &model)
			}
		}
	}

	if err == nil && !model.Blocked {
		mgr := s.Provider.NodeManager()

		if mgr != nil {
			err = setNodeOnline(mgr, &model)
		}
	}
	return model, err
}

func setNodeOnline(mgr node.NodeManager, model *models.Node) error {
	nodeId, err := base64.StdEncoding.DecodeString(model.NodeID)
	if err == nil {
		count := mgr.Count(nodeId)
		model.Online = count > 0
	}
	return err
}

func closeNode(mgr node.NodeManager, model *models.Node) error {
	nodeId, err := base64.StdEncoding.DecodeString(model.NodeID)
	mgr.TraverseNode(nodeId, traverseCloseNode)
	return err
}

func traverseCloseNode(item node.Node) bool {
	item.Close()
	return true
}
