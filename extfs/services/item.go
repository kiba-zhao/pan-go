package services

import (
	"encoding/base64"
	"strconv"

	"pan/app/constant"
	"pan/app/node"

	appModels "pan/app/models"
	appServices "pan/app/services"
	"pan/extfs/models"
	"slices"
	"strings"
)

type ItemService struct {
	Provider                      appServices.NodeManagerProvider
	SettingsExternalService       appServices.SettingsExternalService
	NodeExternalService           appServices.NodeExternalService
	NodeItemInternalService       NodeItemInternalService
	RemoteNodeItemInternalService RemoteNodeItemInternalService
}

const (
	ItemIDSep              = "_"
	ItemTypeNode           = "N"
	ItemTypeRemoteNode     = "RN"
	ItemTypeNodeItem       = "I"
	ItemTypeRemoteNodeItem = "RI"
	ItemTypeFolder         = "D"
	ItemTypeFile           = "F"
)

func (s *ItemService) Search(conditions models.ItemSearchCondition) (total int64, items []models.Item, err error) {
	total = -1
	if conditions.ParentID != nil {
		itemType, linkID := parseItemID(*conditions.ParentID)

		switch itemType {
		case ItemTypeNode:
			nodeItems, err := s.SearchNodeItemWithParentID(conditions.ParentID)
			if err == nil {
				items = append(items, nodeItems...)
			}
		case ItemTypeRemoteNode:
			remoteItems, err := s.SearchRemoteNodeItemByIdWithParentID(linkID, conditions.ParentID)
			if err == nil {
				items = append(items, remoteItems...)
			}
		}
		return
	}

	if len(conditions.ItemTypes) > 0 {
		if slices.Contains(conditions.ItemTypes, ItemTypeNode) {
			items = append(items, s.SelectNode())
		}
		if slices.Contains(conditions.ItemTypes, ItemTypeRemoteNode) {
			remoteItems, remoteNodeErr := s.SearchRemoteNode()
			if remoteNodeErr == nil {
				items = append(items, remoteItems...)
			}
			err = remoteNodeErr
		}
	}
	return
}

func (s *ItemService) SelectNode() models.Item {
	settings := s.SettingsExternalService.Load()
	item := models.Item{}

	item.ID = generateItemID(ItemTypeNode, settings.NodeID)
	item.Name = settings.Name
	item.ItemType = ItemTypeNode
	item.TagQuantity = 0
	item.PendingTagQuantity = 0
	item.Available = true
	return item
}

func (s *ItemService) SearchRemoteNode() ([]models.Item, error) {
	mgr := s.Provider.NodeManager()
	if mgr == nil {
		return nil, nil
	}

	nodeIds := make([]string, 0)
	err := mgr.TraverseNodeID(func(nodeId node.NodeID) error {
		nodeIdBase64 := base64.StdEncoding.EncodeToString(nodeId)
		idx, _ := slices.BinarySearch(nodeIds, nodeIdBase64)
		nodeIds = slices.Insert(nodeIds, idx, nodeIdBase64)
		return nil
	})

	if err != nil || len(nodeIds) <= 0 {
		return nil, err
	}

	items := make([]models.Item, 0)
	err = s.NodeExternalService.TraverseWithNodeIDs(func(model appModels.Node) error {
		idx, ok := slices.BinarySearch(nodeIds, model.NodeID)
		if !ok {
			return constant.ErrInternalError
		}
		nodeIds = slices.Delete(nodeIds, idx, idx+1)

		var item models.Item
		item.ID = generateItemID(ItemTypeRemoteNode, model.NodeID)
		item.Name = model.Name
		item.ItemType = ItemTypeRemoteNode
		item.Available = true
		item.TagQuantity = 0
		item.PendingTagQuantity = 0
		item.LinkID = &model.NodeID

		// TODO: set tag quantity

		items = append(items, item)

		return nil
	}, nodeIds)

	if err == nil && len(nodeIds) > 0 {
		for _, nodeId := range nodeIds {
			var item models.Item
			item.ID = generateItemID(ItemTypeRemoteNode, nodeId)
			item.Name = ""
			item.ItemType = ItemTypeRemoteNode
			item.Available = true
			item.TagQuantity = 0
			item.PendingTagQuantity = 0
			item.LinkID = &nodeId

			// TODO: set tag quantity

			items = append(items, item)
		}
	}
	return items, err
}

func (s *ItemService) SearchNodeItemWithParentID(parentId *string) ([]models.Item, error) {
	items := make([]models.Item, 0)

	err := s.NodeItemInternalService.TraverseAll(func(nodeItem models.NodeItem) error {
		linkID := strconv.FormatUint(uint64(nodeItem.ID), 10)
		var item models.Item
		switch nodeItem.FileType {
		case FileTypeFolder:
			item.ItemType = ItemTypeFolder
		case FileTypeFile:
			item.ItemType = ItemTypeFile
		}
		item.ID = generateItemID(ItemTypeNodeItem, linkID)
		item.Name = nodeItem.Name
		item.Available = nodeItem.Available
		item.TagQuantity = 0
		item.PendingTagQuantity = 0
		item.LinkID = &linkID
		item.ParentID = parentId

		// TODO: set tag quantity

		items = append(items, item)
		return nil
	})
	return items, err
}

func (s *ItemService) SearchRemoteNodeItemByIdWithParentID(parentLinkId string, parentId *string) ([]models.Item, error) {

	nodeId, err := base64.StdEncoding.DecodeString(parentLinkId)
	if err != nil {
		return nil, err
	}

	items := make([]models.Item, 0)
	err = s.RemoteNodeItemInternalService.TraverseAllWithNodeID(func(nodeItem *models.RemoteNodeItem) error {
		var item models.Item
		switch nodeItem.FileType {
		case FileTypeFolder:
			item.ItemType = ItemTypeFolder
		case FileTypeFile:
			item.ItemType = ItemTypeFile
		}
		linkId := strconv.FormatUint(uint64(nodeItem.ID), 10)
		item.ID = generateItemID(ItemTypeRemoteNodeItem, linkId)
		item.Name = nodeItem.Name
		item.Available = nodeItem.Available
		item.TagQuantity = 0
		item.PendingTagQuantity = 0
		item.LinkID = &linkId
		item.ParentID = parentId

		// TODO: set tag quantity
		items = append(items, item)
		return nil
	}, nodeId)

	return items, err
}

func generateItemID(itemType string, id string) string {
	return strings.Join([]string{itemType, id}, ItemIDSep)
}

func parseItemID(id string) (string, string) {
	ids := strings.Split(id, ItemIDSep)
	if len(ids) != 2 {
		return "", ""
	}

	switch ids[0] {
	case ItemTypeNode:
		return ids[0], ids[1]
	case ItemTypeRemoteNode:
		return ids[0], ids[1]
	}
	return "", ""
}

// func compareRemoteTagReportWithNodeID(report models.NodeTagReport, nodeId string) int {
// 	ret := cmp.Compare(report.NodeID, nodeId)
// 	return ret
// }
