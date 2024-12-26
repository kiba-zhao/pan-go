package remotefile

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"pan/app/peer"
	"time"

	appnode "pan/app/app_node"
	nodefile "pan/extfs/node_file"
)

var ErrRemoteFileInvalidID = errors.New("remotefile.ParseRemoteFileID Error: Invalid ID")

type RemoteFileService struct {
	RemoteFileBroker *RemoteFileBroker
	NodeFileService  nodefile.NodeFileInternalService
}

func (s *RemoteFileService) Select(id string) (RemoteFile, error) {
	// TODO: implement
	return RemoteFile{}, nil
	//
}

func (s *RemoteFileService) Search(condition RemoteFileSearchCondition) (total int64, items []RemoteFile, err error) {

	peerId, err := appnode.DecodePeerID(condition.PeerID)
	if err != nil {
		return
	}

	recordSearch := RemoteFileRecordSearchCondition{
		ItemID:     uint32(condition.ItemID),
		ParentPath: condition.ParentPath,
	}

	err = s.TraverseRecordWithPeerID(func(record *RemoteFileRecord) error {
		var item RemoteFile

		item.PeerID = condition.PeerID
		item.ItemID = uint(record.ItemID)
		item.Name = record.Name
		item.FileType = record.FileType
		item.FilePath = record.FilePath
		item.ParentPath = record.ParentPath
		item.Size = record.Size
		item.Available = record.Available
		item.CreatedAt = time.Unix(record.CreatedAt, 0)
		item.UpdatedAt = time.Unix(record.UpdatedAt, 0)
		item.ID = GenerateRemoteFileID(peerId, item.ItemID, item.FilePath)

		items = append(items, item)
		return nil
	}, peerId, &recordSearch)

	total = int64(len(items))
	return
}

func (s *RemoteFileService) SearchForTopic(condition *RemoteFileRecordSearchCondition) (*RemoteFileRecordList, error) {
	var condition_ nodefile.NodeFileSearchCondition
	condition_.ItemID = uint(condition.ItemID)
	condition_.ParentPath = condition.ParentPath

	var recordList RemoteFileRecordList
	err := s.NodeFileService.TraverseWithCondition(func(item nodefile.NodeFile) error {
		var record RemoteFileRecord

		record.ID = item.ID
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

func (s *RemoteFileService) SelectForTopic(condition *RemoteFileRecordSelectCondition) (*RemoteFileRecord, error) {
	var condition_ nodefile.NodeFileSelectCondition

	condition_.ItemID = uint(condition.ItemID)
	condition_.ParentPath = condition.ParentPath
	condition_.Name = condition.Name
	fileItem, err := s.NodeFileService.SelectWithCondition(condition_)
	if err != nil {
		return nil, err
	}

	var record RemoteFileRecord
	record.ID = fileItem.ID
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

func (s *RemoteFileService) TraverseRecordWithPeerID(traverseFn func(record *RemoteFileRecord) error, peerId peer.PeerID, condition *RemoteFileRecordSearchCondition) error {

	recordList, err := s.RemoteFileBroker.Search(peerId, condition)
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

func (s *RemoteFileService) SelectWithCondition(peerId peer.PeerID, condition *RemoteFileRecordSelectCondition) (*RemoteFileRecord, error) {
	return s.RemoteFileBroker.Select(peerId, condition)
}

func GenerateRemoteFileID(peerId peer.PeerID, itemID uint, filePath string) string {
	filePathBytes := []byte(filePath)
	idBytes := make([]byte, 6+len(peerId)+len(filePathBytes))

	// set peerId length + peerId
	peerIdLen := len(peerId)
	binary.BigEndian.PutUint16(idBytes, uint16(peerIdLen))
	copy(idBytes[2:], peerId)

	// set itemID
	binary.BigEndian.PutUint32(idBytes[len(peerId)+2:], uint32(itemID))
	// set filePath
	copy(idBytes[len(peerId)+6:], filePathBytes)
	return base64.StdEncoding.EncodeToString(idBytes)
}

func ParseRemoteFileID(id string) (peer.PeerID, uint, string, error) {
	idBytes, err := base64.StdEncoding.DecodeString(id)
	if err == nil && len(idBytes) < 6 {
		err = ErrRemoteFileInvalidID
	}
	if err != nil {
		return nil, 0, "", err
	}
	// parse peerId
	peerIdLen := binary.BigEndian.Uint16(idBytes)
	peerId := idBytes[2 : peerIdLen+2]
	// parse itemID
	itemId := binary.BigEndian.Uint32(idBytes[peerIdLen+2:])
	// parse filePath
	filePath := string(idBytes[peerIdLen+6:])
	return peerId, uint(itemId), filePath, err
}
