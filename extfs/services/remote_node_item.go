package services

import (
	"bytes"
	"encoding/base64"
	"io"
	appConstant "pan/app/constant"

	appNode "pan/app/node"
	"pan/extfs/constant"
	"pan/extfs/models"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

type RemoteNodeItemService struct {
	NodeModule      appNode.NodeModule
	NodeItemService NodeItemInternalService
	NodeScopeModule appNode.NodeScopeModule
}

func (s *RemoteNodeItemService) Search(condition models.RemoteNodeItemSearchCondition) (total int64, items []models.RemoteNodeItem, err error) {
	nodeId, err := base64.StdEncoding.DecodeString(condition.NodeID)
	if err != nil {
		return
	}

	err = s.TraverseRecordWithNodeID(func(record *models.RemoteNodeItemRecord) error {
		var item models.RemoteNodeItem

		item.NodeID = condition.NodeID
		item.ItemID = uint(record.ID)
		item.Name = record.Name
		item.FileType = record.FileType
		item.Size = record.Size
		item.Available = record.Available
		item.CreatedAt = time.Unix(record.CreatedAt, 0)
		item.UpdatedAt = time.Unix(record.UpdatedAt, 0)

		item.ID = generateRemoteNodeItemId(item.NodeID, item.ItemID)
		items = append(items, item)
		return nil
	}, nodeId)

	if err != nil {
		return
	}

	total = int64(len(items))
	return
}

func (s *RemoteNodeItemService) SelectAllForNode() (models.RemoteNodeItemRecordList, error) {

	var recordList models.RemoteNodeItemRecordList

	err := s.NodeItemService.TraverseAll(func(nodeItem models.NodeItem) error {
		var record models.RemoteNodeItemRecord
		record.ID = int32(nodeItem.ID)
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

func (s *RemoteNodeItemService) SelectForNode(condition *models.RemoteNodeItemRecordSelectCondition) (*models.RemoteNodeItemRecord, error) {
	var nodeItem models.NodeItem
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
	var record models.RemoteNodeItemRecord
	record.ID = int32(nodeItem.ID)
	record.Name = nodeItem.Name
	record.FileType = nodeItem.FileType
	record.Size = nodeItem.Size
	record.Available = nodeItem.Available
	record.CreatedAt = nodeItem.CreatedAt.Unix()
	record.UpdatedAt = nodeItem.UpdatedAt.Unix()
	return &record, err
}

var RequestAllRemoteItems = []byte("select_all_remote_items")

func (s *RemoteNodeItemService) TraverseRecordWithNodeID(traverseFn func(record *models.RemoteNodeItemRecord) error, nodeId appNode.NodeID) error {
	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestAllRemoteItems)
	request := appNode.NewRequest(requestName, nil)

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

	var remoteNodeItemRecordList models.RemoteNodeItemRecordList
	err = proto.Unmarshal(data, &remoteNodeItemRecordList)
	if err != nil {
		return err
	}

	for _, item := range remoteNodeItemRecordList.Items {
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}

var RequestRemoteItem = []byte("select_remote_item")

func (s *RemoteNodeItemService) SelectWithCondition(nodeId appNode.NodeID, condition *models.RemoteNodeItemRecordSelectCondition) (*models.RemoteNodeItemRecord, error) {
	requestBytes, err := proto.Marshal(condition)
	if err != nil {
		return nil, err
	}

	scope := s.NodeScopeModule.NodeScope()
	requestName := appNode.GenerateRouteName(scope, RequestRemoteItem)
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

	var record models.RemoteNodeItemRecord
	err = proto.Unmarshal(data, &record)

	return &record, err
}

func generateRemoteNodeItemId(nodeId string, itemId uint) string {
	idStr := strings.Join([]string{nodeId, strconv.FormatUint(uint64(itemId), 10)}, constant.FileItemSep)
	return idStr
}
