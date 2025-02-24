package appnode

import (
	"encoding/base64"
	"errors"
	"pan/app/peer"
	"strings"
)

var ErrAppNodeBlocked = errors.New("appnode.AppNodeService Error: App Node Blocked")

type AppNodeExternalService interface {
	TraverseAll(func(AppNode) error) error
	SelectByName(string) (AppNode, error)
}

type AppNodeService struct {
	AppNodeRepo AppNodeRepository
	PeerModule  peer.PeerModule
}

func (s *AppNodeService) Search(conditions AppNodeSearchCondition) (total int64, items []AppNode, err error) {
	total, items, err = s.AppNodeRepo.Search(conditions)
	if err != nil {
		return
	}

	if conditions.Blocked != nil && *conditions.Blocked {
		return
	}

	peerModule := s.PeerModule
	if peerModule == nil {
		return
	}
	items_ := make([]AppNode, 0)
	for _, item := range items {
		if !item.Blocked {
			setPeerOnline(peerModule, &item)
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
		peerModule := s.PeerModule
		if peerModule != nil {
			err = setPeerOnline(peerModule, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) SelectByName(name string) (AppNode, error) {
	model, err := s.AppNodeRepo.SelectByName(name)
	if err == nil && !model.Blocked {
		peerModule := s.PeerModule
		if peerModule != nil {
			err = setPeerOnline(peerModule, &model)
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
		peerModule := s.PeerModule
		if peerModule != nil {
			err = purgeWithPeerID(peerModule, &model)
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
			peerModule := s.PeerModule
			if peerModule != nil {
				err = purgeWithPeerID(peerModule, &model)
			}
		}
	}

	if err == nil && !model.Blocked {
		peerModule := s.PeerModule
		if peerModule != nil {
			err = setPeerOnline(peerModule, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) AccessWithPeerID(peerId peer.PeerID) error {
	peerId_ := EncodePeerID(peerId)
	node_, err := s.AppNodeRepo.SelectByPeerID(peerId_)
	if err == nil && node_.Blocked {
		err = ErrAppNodeBlocked
	}
	return err
}

func (s *AppNodeService) TraverseAll(traverseFn func(model AppNode) error) error {
	return s.AppNodeRepo.TraverseAll(func(an AppNode) error {
		var err error
		peerModule := s.PeerModule
		if peerModule != nil {
			err = setPeerOnline(peerModule, &an)
		}
		if err != nil {
			return err
		}
		return traverseFn(an)
	})
}

func setPeerOnline(peerModule peer.PeerModule, model *AppNode) error {
	peerId, err := DecodePeerID(model.PeerID)
	if err == nil {
		model.Online = peerModule.CanReach(peerId)
	}
	return err
}

func purgeWithPeerID(peerModule peer.PeerModule, model *AppNode) error {
	peerId, err := DecodePeerID(model.PeerID)
	if err == nil {
		err = peerModule.Purge(peerId)
	}
	return err
}

func EncodePeerID(peerId peer.PeerID) string {
	return base64.RawURLEncoding.EncodeToString(peerId)
}

func DecodePeerID(peerId string) (peer.PeerID, error) {
	return base64.RawURLEncoding.DecodeString(peerId)
}
