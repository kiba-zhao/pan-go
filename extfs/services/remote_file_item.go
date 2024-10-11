package services

import (
	"bytes"
	"encoding/base64"
	"io"
	appConstant "pan/app/constant"
	appNode "pan/app/node"
	"pan/extfs/models"
	"time"

	"google.golang.org/protobuf/proto"
)

type RemoteFileItemService struct {
	NodeModule      appNode.NodeModule
	FileItemService FileItemInternalService
}

func (s *RemoteFileItemService) Search(condition models.RemoteFileItemSearchCondition) (total int64, items []models.RemoteFileItem, err error) {

	nodeId, err := base64.StdEncoding.DecodeString(condition.NodeID)
	if err != nil {
		return
	}

	recordSearch := models.RemoteFileItemRecordSearchCondition{
		ItemID:     int32(condition.ItemID),
		ParentPath: condition.ParentPath,
	}

	err = s.TraverseRecordWithNodeID(func(record *models.RemoteFileItemRecord) error {
		var item models.RemoteFileItem

		item.NodeID = condition.NodeID
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
	}, nodeId, &recordSearch)

	total = int64(len(items))
	return
}

func (s *RemoteFileItemService) SearchForNode(condition *models.RemoteFileItemRecordSearchCondition) (*models.RemoteFileItemRecordList, error) {
	var condition_ models.FileItemSearchCondition
	condition_.ItemID = uint(condition.ItemID)
	condition_.ParentPath = condition.ParentPath

	var recordList models.RemoteFileItemRecordList
	err := s.FileItemService.TraverseWithCondition(func(item models.FileItem) error {
		var record models.RemoteFileItemRecord

		record.ID = item.ID
		record.Name = item.Name
		record.FilePath = item.FilePath
		record.ParentPath = item.ParentPath
		record.Size = item.Size
		record.FileType = item.FileType
		record.ItemID = int32(item.ItemID)
		record.Available = item.Available
		record.CreatedAt = item.CreatedAt.Unix()
		record.UpdatedAt = item.UpdatedAt.Unix()
		recordList.Items = append(recordList.Items, &record)
		return nil
	}, condition_)

	return &recordList, err
}

var RequestAllRemoteFileItems = []byte("extfs/select_all_remote_file_items")

func (s *RemoteFileItemService) TraverseRecordWithNodeID(traverseFn func(record *models.RemoteFileItemRecord) error, nodeId appNode.NodeID, condition *models.RemoteFileItemRecordSearchCondition) error {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return err
	}

	request := appNode.NewRequest(RequestAllRemoteFileItems, bytes.NewReader(requestBytes))
	response, err := s.NodeModule.Do(nodeId, request)

	if err != nil {
		return err
	}
	if response.Code() != appConstant.CodeOK {
		return appConstant.ErrInternalError
	}
	data, err := io.ReadAll(response.Body())
	if err != nil {
		return err
	}

	var recordList models.RemoteFileItemRecordList
	err = proto.Unmarshal(data, &recordList)
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
