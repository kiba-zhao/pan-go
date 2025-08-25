package node

import (
	"errors"
	"pan/lib/peer"
	"slices"
	"strings"
)

var ErrAppNodeBlocked = errors.New("appnode.AppNodeService Error: App Node Blocked")

type AppNodeExternalService interface {
	TraverseAll(func(AppNode) error) error
	SelectByName(string) (AppNode, error)
}

type AppNodeService struct {
	AppNodeRepo AppNodeRepository
	PeerClient  peer.PeerClient
}

func (s *AppNodeService) Search(conditions AppNodeSearchCondition) (total int64, items []AppNode, err error) {
	total, items, err = s.AppNodeRepo.Search(conditions)
	if err != nil {
		return
	}

	if conditions.Blocked != nil && *conditions.Blocked {
		return
	}

	peerClient := s.PeerClient
	if peerClient == nil {
		return
	}
	items_ := make([]AppNode, 0)
	for _, item := range items {
		if !item.Blocked {
			setPeerOnline(peerClient, &item)
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
	if err != nil {
		return model, err
	}
	for _, networkAddr := range model.NetworkAddrs {
		model.NetworkAddrTexts = append(model.NetworkAddrTexts, networkAddr.Address)
	}

	peerClient := s.PeerClient
	if !model.Blocked && peerClient != nil {
		err = setPeerOnline(peerClient, &model)
	}
	return model, err
}

func (s *AppNodeService) SelectByName(name string) (AppNode, error) {
	model, err := s.AppNodeRepo.SelectByName(name)
	if err == nil && !model.Blocked {
		peerClient := s.PeerClient
		if peerClient != nil {
			err = setPeerOnline(peerClient, &model)
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
		peerClient := s.PeerClient
		if peerClient != nil {
			err = purgeWithPeerID(peerClient, &model)
		}
	}
	return err
}

func (s *AppNodeService) Create(fields AppNodeFields) (AppNode, error) {

	var model AppNode
	model.Name = fields.Name
	model.PeerID = fields.PeerID
	model.Blocked = fields.Blocked != nil && *fields.Blocked
	for _, addrText := range fields.NetworkAddrs {
		addrText = strings.Trim(addrText, " ")
		if len(addrText) > 0 {
			networkAddr := NetworkAddr{
				Address: addrText,
			}
			model.NetworkAddrs = append(model.NetworkAddrs, networkAddr)
		}
	}

	model, err := s.AppNodeRepo.Create(model)
	if err != nil {
		return model, err
	}
	model.NetworkAddrTexts = fields.NetworkAddrs
	return model, err

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

	if fields.NetworkAddrs != nil {
		dirty = true
		addrs := make([]NetworkAddr, 0)
		fieldAddrs := slices.Clone(fields.NetworkAddrs)
		for _, networkAddr := range model.NetworkAddrs {
			matchIdx := -1
			for idx, addrText := range fieldAddrs {
				addrText = strings.Trim(addrText, " ")
				if len(addrText) > 0 && strings.Compare(networkAddr.Address, addrText) == 0 {
					matchIdx = idx
					break
				}
			}

			if matchIdx >= 0 {
				addrs = append(addrs, networkAddr)
				fieldAddrs = slices.Delete(fieldAddrs, matchIdx, matchIdx+1)
			}
		}
		for _, addrText := range fieldAddrs {
			addrText = strings.Trim(addrText, " ")
			if len(addrText) > 0 {
				networkAddr := NetworkAddr{
					AppNodeID: model.ID,
					Address:   addrText,
				}
				addrs = append(addrs, networkAddr)
			}
		}
		model.NetworkAddrs = addrs
	}

	if dirty {
		model, err = s.AppNodeRepo.Update(model)
		if err == nil && needClosed {
			peerClient := s.PeerClient
			if peerClient != nil {
				err = purgeWithPeerID(peerClient, &model)
			}
		}
		model.NetworkAddrTexts = fields.NetworkAddrs
	}

	if err == nil && !model.Blocked {
		peerClient := s.PeerClient
		if peerClient != nil {
			err = setPeerOnline(peerClient, &model)
		}
	}
	return model, err
}

func (s *AppNodeService) AccessWithPeerID(peerId peer.PeerID) error {
	peerId_ := peer.EncodePeerID(peerId)
	node_, err := s.AppNodeRepo.SelectByPeerID(peerId_)
	if err == nil && node_.Blocked {
		err = ErrAppNodeBlocked
	}
	return err
}

func (s *AppNodeService) TraverseAll(traverseFn func(model AppNode) error) error {
	return s.AppNodeRepo.TraverseAll(func(an AppNode) error {
		var err error
		peerClient := s.PeerClient
		if peerClient != nil {
			err = setPeerOnline(peerClient, &an)
		}
		if err != nil {
			return err
		}
		return traverseFn(an)
	})
}

func setPeerOnline(peerClient peer.PeerClient, model *AppNode) error {
	peerId, err := peer.DecodePeerID(model.PeerID)
	if err == nil {
		model.Online = peerClient.CanReach(peerId)
	}
	return err
}

func purgeWithPeerID(peerClient peer.PeerClient, model *AppNode) error {
	peerId, err := peer.DecodePeerID(model.PeerID)
	if err == nil {
		err = peerClient.Purge(peerId)
	}
	return err
}
