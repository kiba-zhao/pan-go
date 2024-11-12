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
	NodeScopeModule appNode.NodeScopeModule
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

func (s *RemoteFileItemService) SelectForNode(condition *models.RemoteFileItemRecordSelectCondition) (*models.RemoteFileItemRecord, error) {
	var condition_ models.FileItemSelectCondition

	condition_.ItemID = uint(condition.ItemID)
	condition_.ParentPath = condition.ParentPath
	condition_.Name = condition.Name
	fileItem, err := s.FileItemService.SelectWithCondition(condition_)
	if err != nil {
		return nil, err
	}

	var record models.RemoteFileItemRecord
	record.ID = fileItem.ID
	record.Name = fileItem.Name
	record.FilePath = fileItem.FilePath
	record.ParentPath = fileItem.ParentPath
	record.Size = fileItem.Size
	record.FileType = fileItem.FileType
	record.ItemID = int32(fileItem.ItemID)
	record.Available = fileItem.Available
	record.CreatedAt = fileItem.CreatedAt.Unix()
	record.UpdatedAt = fileItem.UpdatedAt.Unix()
	return &record, nil
}

var RequestAllRemoteFileItems = []byte("select_all_remote_file_items")

func (s *RemoteFileItemService) TraverseRecordWithNodeID(traverseFn func(record *models.RemoteFileItemRecord) error, nodeId appNode.NodeID, condition *models.RemoteFileItemRecordSearchCondition) error {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return err
	}

	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestAllRemoteFileItems)
	request := appNode.NewRequest(requestName, bytes.NewReader(requestBytes))
	res, err := s.NodeModule.Do(nodeId, request)
	if err != nil {
		return err
	}
	defer res.Close()

	if res.Code() != appConstant.CodeOK {
		return appConstant.ErrInternalError
	}
	data, err := io.ReadAll(res)
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

var RequestRemoteFileItem = []byte("select_remote_file_item")

func (s *RemoteFileItemService) SelectWithCondition(nodeId appNode.NodeID, condition *models.RemoteFileItemRecordSelectCondition) (*models.RemoteFileItemRecord, error) {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return nil, err
	}

	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestRemoteFileItem)
	request := appNode.NewRequest(requestName, bytes.NewReader(requestBytes))
	res, err := s.NodeModule.Do(nodeId, request)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	if res.Code() != appConstant.CodeOK {
		return nil, appConstant.ErrInternalError
	}
	data, err := io.ReadAll(res)
	if err != nil {
		return nil, err
	}

	var record models.RemoteFileItemRecord
	err = proto.Unmarshal(data, &record)
	return &record, err
}
