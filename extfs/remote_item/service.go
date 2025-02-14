package remoteitem

import (
	"errors"
	appnode "pan/app/app_node"
	"pan/app/peer"

	nodeitem "pan/extfs/node_item"
	"time"
)

var ErrRemoteItemInvalidID = errors.New("remoteitem.ParseRemoteItemId Error: Invalid ID")

type RemoteItemService struct {
	RemoteItemBroker *RemoteItemBroker
	NodeItemService  nodeitem.NodeItemInternalService
}

func (s *RemoteItemService) IsNotExist(err error) bool {
	if err == nil {
		return false
	}

	if peerErr, ok := err.(*peer.PeerError); ok && peerErr.Code() == peer.CodeNotFound {
		return true
	}
	return s.NodeItemService.IsNotExist(err)
}

func (s *RemoteItemService) Select(peerId string, itemId uint) (RemoteItem, error) {

	peerIdBytes, err := appnode.DecodePeerID(peerId)
	if err != nil {
		return RemoteItem{}, err
	}
	itemId32 := uint32(itemId)
	record, err := s.SelectWithCondition(peerIdBytes, &RemoteItemRecordSelectCondition{ID: &itemId32})

	if err != nil {
		return RemoteItem{}, err
	}

	var item RemoteItem

	item.PeerID = peerId
	item.ItemID = uint(record.ID)
	item.Name = record.Name
	item.FileType = record.FileType
	item.Size = record.Size
	item.Available = record.Available
	item.CreatedAt = time.Unix(record.CreatedAt, 0)
	item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

	return item, nil
}

func (s *RemoteItemService) Search(peerId string) (total int64, items []RemoteItem, err error) {
	peerIdBytes, err := appnode.DecodePeerID(peerId)
	if err != nil {
		return
	}

	err = s.TraverseRecordWithPeerID(func(record *RemoteItemRecord) error {
		var item RemoteItem

		item.PeerID = peerId
		item.ItemID = uint(record.ID)
		item.Name = record.Name
		item.FileType = record.FileType
		item.Size = record.Size
		item.Available = record.Available
		item.CreatedAt = time.Unix(record.CreatedAt, 0)
		item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

		items = append(items, item)
		return nil
	}, peerIdBytes)

	if err != nil {
		return
	}

	total = int64(len(items))
	return
}

func (s *RemoteItemService) SelectAllForTopic() (RemoteItemRecordList, error) {

	var recordList RemoteItemRecordList

	err := s.NodeItemService.TraverseAll(func(nodeItem nodeitem.NodeItem) error {
		var record RemoteItemRecord
		record.ID = uint32(nodeItem.ID)
		record.Name = nodeItem.Name
		record.FileType = nodeItem.FileType
		record.Size = nodeItem.Size
		record.Available = nodeItem.Available
		record.CreatedAt = nodeItem.CreatedAt.Unix()
		record.UpdatedAt = nodeItem.UpdatedAt.Unix()

		recordList.Items = append(recordList.Items, &record)
		return nil
	})
	return recordList, err
}

func (s *RemoteItemService) SelectForTopic(condition *RemoteItemRecordSelectCondition) (*RemoteItemRecord, error) {
	var nodeItem nodeitem.NodeItem
	var err error
	if condition.Name != nil {
		nodeItem, err = s.NodeItemService.SelectByName(*condition.Name)
	}
	if condition.ID != nil {
		nodeItem, err = s.NodeItemService.Select(uint(*condition.ID))
	}

	if err != nil {
		return nil, err
	}
	var record RemoteItemRecord
	record.ID = uint32(nodeItem.ID)
	record.Name = nodeItem.Name
	record.FileType = nodeItem.FileType
	record.Size = nodeItem.Size
	record.Available = nodeItem.Available
	record.CreatedAt = nodeItem.CreatedAt.Unix()
	record.UpdatedAt = nodeItem.UpdatedAt.Unix()
	return &record, err
}

func (s *RemoteItemService) TraverseRecordWithPeerID(traverseFn func(record *RemoteItemRecord) error, peerId peer.PeerID) error {

	remoteItemRecordList, err := s.RemoteItemBroker.SelectAll(peerId)
	if err != nil {
		return err
	}
	for _, item := range remoteItemRecordList.Items {
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *RemoteItemService) SelectWithCondition(peerId peer.PeerID, condition *RemoteItemRecordSelectCondition) (*RemoteItemRecord, error) {
	return s.RemoteItemBroker.Select(peerId, condition)
}
