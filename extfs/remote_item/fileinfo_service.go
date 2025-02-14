package remoteitem

import (
	"errors"
	appnode "pan/app/app_node"
	"pan/app/peer"
	"time"

	nodeitem "pan/extfs/node_item"
)

var ErrRemoteFileInfoInvalidID = errors.New("remoteitem.ParseRemoteFileInfoID Error: Invalid ID")

type RemoteFileInfoInternalService interface {
	IsNotExist(err error) bool
	SelectWithCondition(peerId peer.PeerID, condition *RemoteFileInfoRecordSelectCondition) (*RemoteFileInfoRecord, error)
}

type RemoteFileInfoService struct {
	RemoteFileInfoBroker *RemoteFileInfoBroker
	NodeFileInfoService  nodeitem.NodeFileInfoInternalService
}

func (s *RemoteFileInfoService) IsNotExist(err error) bool {
	if peerErr, ok := err.(*peer.PeerError); ok && peerErr.Code() == peer.CodeNotFound {
		return true
	}
	return s.NodeFileInfoService.IsNotExist(err)
}

func (s *RemoteFileInfoService) Select(peerId string, itemId uint, filePath string) (RemoteFileInfo, error) {

	peerIdBytes, err := appnode.DecodePeerID(peerId)
	if err != nil {
		return RemoteFileInfo{}, err
	}

	var condition RemoteFileInfoRecordSelectCondition

	condition.ItemID = uint32(itemId)
	condition.FilePath = filePath

	record, err := s.SelectWithCondition(peerIdBytes, &condition)
	if err != nil {
		return RemoteFileInfo{}, err
	}

	var item RemoteFileInfo

	item.PeerID = peerId
	item.ItemID = uint(record.ItemID)
	item.Name = record.Name
	item.FileType = record.FileType
	item.FilePath = record.FilePath
	item.ParentPath = record.ParentPath
	item.Size = record.Size
	item.Available = record.Available
	item.CreatedAt = time.Unix(record.CreatedAt, 0)
	item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

	return item, nil

}

func (s *RemoteFileInfoService) Search(peerId string, itemId uint, parentPath string) (total int64, items []RemoteFileInfo, err error) {

	peerIdBytes, err := appnode.DecodePeerID(peerId)
	if err != nil {
		return
	}

	recordSearch := RemoteFileInfoRecordSearchCondition{
		ItemID:     uint32(itemId),
		ParentPath: parentPath,
	}

	err = s.TraverseRecordWithPeerID(func(record *RemoteFileInfoRecord) error {
		var item RemoteFileInfo

		item.PeerID = peerId
		item.ItemID = uint(record.ItemID)
		item.Name = record.Name
		item.FileType = record.FileType
		item.FilePath = record.FilePath
		item.ParentPath = record.ParentPath
		item.Size = record.Size
		item.Available = record.Available
		item.CreatedAt = time.Unix(record.CreatedAt, 0)
		item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

		items = append(items, item)
		return nil
	}, peerIdBytes, &recordSearch)

	total = int64(len(items))
	return
}

func (s *RemoteFileInfoService) SearchForTopic(condition *RemoteFileInfoRecordSearchCondition) (*RemoteFileInfoRecordList, error) {
	var condition_ nodeitem.NodeFileInfoSearchCondition
	condition_.ItemID = uint(condition.ItemID)
	condition_.ParentPath = condition.ParentPath

	var recordList RemoteFileInfoRecordList
	err := s.NodeFileInfoService.TraverseWithCondition(func(item nodeitem.NodeFileInfo) error {
		var record RemoteFileInfoRecord

		record.Name = item.Name
		record.FilePath = item.FilePath
		record.ParentPath = item.ParentPath
		record.Size = item.Size
		record.FileType = item.FileType
		record.ItemID = uint32(item.ItemID)
		record.Available = item.Available
		record.CreatedAt = item.CreatedAt.Unix()
		record.UpdatedAt = item.UpdatedAt.Unix()

		recordList.Items = append(recordList.Items, &record)
		return nil
	}, condition_)

	return &recordList, err
}

func (s *RemoteFileInfoService) SelectForTopic(condition *RemoteFileInfoRecordSelectCondition) (*RemoteFileInfoRecord, error) {

	fileItem, err := s.NodeFileInfoService.Select(uint(condition.ItemID), condition.FilePath)
	if err != nil {
		return nil, err
	}

	var record RemoteFileInfoRecord

	record.Name = fileItem.Name
	record.FilePath = fileItem.FilePath
	record.ParentPath = fileItem.ParentPath
	record.Size = fileItem.Size
	record.FileType = fileItem.FileType
	record.ItemID = uint32(fileItem.ItemID)
	record.Available = fileItem.Available
	record.CreatedAt = fileItem.CreatedAt.Unix()
	record.UpdatedAt = fileItem.UpdatedAt.Unix()
	return &record, nil
}

func (s *RemoteFileInfoService) TraverseRecordWithPeerID(traverseFn func(record *RemoteFileInfoRecord) error, peerId peer.PeerID, condition *RemoteFileInfoRecordSearchCondition) error {

	recordList, err := s.RemoteFileInfoBroker.Search(peerId, condition)
	if err != nil {
		return err
	}

	for _, item := range recordList.Items {
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *RemoteFileInfoService) SelectWithCondition(peerId peer.PeerID, condition *RemoteFileInfoRecordSelectCondition) (*RemoteFileInfoRecord, error) {
	return s.RemoteFileInfoBroker.Select(peerId, condition)
}
