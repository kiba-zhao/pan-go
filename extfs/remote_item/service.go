package remoteitem

import (
	"encoding/base64"
	"encoding/binary"
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

func (s *RemoteItemService) Select(id string) (RemoteItem, error) {

	itemId, peerId, err := ParseRemoteItemId(id)
	if err != nil {
		return RemoteItem{}, err
	}
	itemId32 := uint32(itemId)
	record, err := s.SelectWithCondition(peerId, &RemoteItemRecordSelectCondition{ID: &itemId32})

	if err != nil {
		return RemoteItem{}, err
	}

	var item RemoteItem

	item.ID = id
	item.PeerID = appnode.EncodePeerID(peerId)
	item.ItemID = uint(record.ID)
	item.Name = record.Name
	item.FileType = record.FileType
	item.Size = record.Size
	item.Available = record.Available
	item.CreatedAt = time.Unix(record.CreatedAt, 0)
	item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

	return item, nil
}

func (s *RemoteItemService) Search(condition RemoteItemSearchCondition) (total int64, items []RemoteItem, err error) {
	peerId, err := appnode.DecodePeerID(condition.PeerID)
	if err != nil {
		return
	}

	err = s.TraverseRecordWithPeerID(func(record *RemoteItemRecord) error {
		var item RemoteItem

		item.PeerID = condition.PeerID
		item.ItemID = uint(record.ID)
		item.Name = record.Name
		item.FileType = record.FileType
		item.Size = record.Size
		item.Available = record.Available
		item.CreatedAt = time.Unix(record.CreatedAt, 0)
		item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

		item.ID = GenerateRemoteItemId(peerId, item.ItemID)
		items = append(items, item)
		return nil
	}, peerId)

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

func GenerateRemoteItemId(peerId peer.PeerID, itemId uint) string {
	idBytes := make([]byte, 4+len(peerId))
	binary.BigEndian.PutUint32(idBytes, uint32(itemId))
	copy(idBytes[4:], peerId)
	return base64.StdEncoding.EncodeToString(idBytes)
}

func ParseRemoteItemId(id string) (uint, peer.PeerID, error) {
	idBytes, err := base64.StdEncoding.DecodeString(id)
	if err == nil && len(idBytes) < 4 {
		err = ErrRemoteItemInvalidID
	}
	if err != nil {
		return 0, nil, err
	}
	itemId := binary.BigEndian.Uint32(idBytes)
	peerId := idBytes[4:]
	return uint(itemId), peerId, err
}
