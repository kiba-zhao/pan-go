package services

import (
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

var RequestAllRemoteItems = []byte("extfs/select_all_remote_items")

func (s *RemoteNodeItemService) TraverseRecordWithNodeID(traverseFn func(record *models.RemoteNodeItemRecord) error, nodeId appNode.NodeID) error {
	request := appNode.NewRequest(RequestAllRemoteItems, nil)

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

func generateRemoteNodeItemId(nodeId string, itemId uint) string {
	idStr := strings.Join([]string{nodeId, strconv.FormatUint(uint64(itemId), 10)}, constant.FileItemSep)
	return idStr
}
