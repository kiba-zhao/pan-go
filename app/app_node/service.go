package appnode

import (
	"encoding/base64"
	"errors"
	"pan/app/peer"
	"strings"
)

var ErrAppNodeBlocked = errors.New("appnode.AppNodeService Error: App Node Blocked")

type PeerManagerProvider interface {
	PeerManager() peer.PeerManager
}

type AppNodeExternalService interface {
	TraverseWithPeerIDs(func(AppNode) error, []string) error
	SelectByName(string) (AppNode, error)
}

type AppNodeService struct {
	AppNodeRepo AppNodeRepository
	Provider    PeerManagerProvider
}

func (s *AppNodeService) Search(conditions AppNodeSearchCondition) (total int64, items []AppNode, err error) {
	total, items, err = s.AppNodeRepo.Search(conditions)
	if err != nil {
		return
	}

	if conditions.Blocked != nil && *conditions.Blocked {
		return
	}

	mgr := s.Provider.PeerManager()
	if mgr == nil {
		return
	}
	items_ := make([]AppNode, 0)
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

func (s *AppNodeService) Select(id uint) (AppNode, error) {

	model, err := s.AppNodeRepo.Select(id)
	if err == nil && !model.Blocked {
		mgr := s.Provider.PeerManager()
		if mgr != nil {
			err = setNodeOnline(mgr, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) SelectByName(name string) (AppNode, error) {
	model, err := s.AppNodeRepo.SelectByName(name)
	if err == nil && !model.Blocked {
		mgr := s.Provider.PeerManager()
		if mgr != nil {
			err = setNodeOnline(mgr, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) Delete(id uint) error {
	model, err := s.AppNodeRepo.Select(id)
	if err != nil {
		return err
	}
	err = s.AppNodeRepo.Delete(model)
	if err != nil {
		return err
	}
	if !model.Blocked {
		mgr := s.Provider.PeerManager()
		if mgr != nil {
			err = closeNode(mgr, &model)
		}
	}
	return err
}

func (s *AppNodeService) Create(fields AppNodeFields) (AppNode, error) {

	var model AppNode
	model.Name = fields.Name
	model.PeerID = fields.PeerID
	model.Blocked = fields.Blocked != nil && *fields.Blocked

	return s.AppNodeRepo.Save(model)

}

func (s *AppNodeService) Update(id uint, fields AppNodeFields) (AppNode, error) {

	model, err := s.AppNodeRepo.Select(id)
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
		model, err = s.AppNodeRepo.Save(model)
		if err == nil && needClosed {
			mgr := s.Provider.PeerManager()
			if mgr != nil {
				err = closeNode(mgr, &model)
			}
		}
	}

	if err == nil && !model.Blocked {
		mgr := s.Provider.PeerManager()

		if mgr != nil {
			err = setNodeOnline(mgr, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) AccessWithPeerID(peerId peer.PeerID) error {
	peerId_ := base64.StdEncoding.EncodeToString(peerId)
	node_, err := s.AppNodeRepo.SelectByPeerID(peerId_)
	if err == nil && node_.Blocked {
		err = ErrAppNodeBlocked
	}
	return err
}

func (s *AppNodeService) TraverseWithPeerIDs(traverseFn func(model AppNode) error, peerIds []string) error {
	return s.AppNodeRepo.TraverseWithPeerIDs(traverseFn, peerIds)
}

func setNodeOnline(mgr peer.PeerManager, model *AppNode) error {
	peerId, err := base64.StdEncoding.DecodeString(model.PeerID)
	if err == nil {
		count := mgr.Count(peerId)
		model.Online = count > 0
	}
	return err
}

func closeNode(mgr peer.PeerManager, model *AppNode) error {
	peerId, err := base64.StdEncoding.DecodeString(model.PeerID)
	mgr.TraversePeerNode(peerId, traverseCloseAppNode)
	return err
}

func traverseCloseAppNode(item peer.PeerNode) bool {
	item.Close()
	return true
}
