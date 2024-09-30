package services

import (
	"encoding/base64"
	"pan/app/constant"
	appModels "pan/app/models"
	"pan/app/node"
	appServices "pan/app/services"
	"pan/extfs/models"
	"slices"
)

type RemoteNodeService struct {
	Provider            appServices.NodeManagerProvider
	NodeExternalService appServices.NodeExternalService
}

func (s *RemoteNodeService) SelectAll() (int64, []models.RemoteNode, error) {
	mgr := s.Provider.NodeManager()
	if mgr == nil {
		return 0, nil, constant.ErrUnavailable
	}

	nodeIds := make([]string, 0)
	err := mgr.TraverseNodeID(func(nodeId node.NodeID) error {
		nodeIdBase64 := base64.StdEncoding.EncodeToString(nodeId)
		idx, _ := slices.BinarySearch(nodeIds, nodeIdBase64)
		nodeIds = slices.Insert(nodeIds, idx, nodeIdBase64)
		return nil
	})

	if err != nil || len(nodeIds) <= 0 {
		return 0, nil, err
	}

	remotes := make([]models.RemoteNode, 0)
	err = s.NodeExternalService.TraverseWithNodeIDs(func(model appModels.Node) error {
		idx, ok := slices.BinarySearch(nodeIds, model.NodeID)
		if !ok {
			return constant.ErrInternalError
		}
		nodeIds = slices.Delete(nodeIds, idx, idx+1)

		var remote models.RemoteNode
		remote.Name = model.Name
		remote.NodeID = model.NodeID
		remote.Available = true
		remote.TagQuantity = 0
		remote.PendingTagQuantity = 0
		remote.CreatedAt = model.CreatedAt
		remote.UpdatedAt = model.UpdatedAt

		// TODO: set tag quantity

		remotes = append(remotes, remote)

		return nil
	}, nodeIds)

	if err == nil && len(nodeIds) > 0 {
		for _, nodeId := range nodeIds {
			var remote models.RemoteNode
			remote.NodeID = nodeId
			remote.Available = true
			remote.TagQuantity = 0
			remote.PendingTagQuantity = 0

			// TODO: set tag quantity

			remotes = append(remotes, remote)
		}
	}

	return 0, remotes, err
}
