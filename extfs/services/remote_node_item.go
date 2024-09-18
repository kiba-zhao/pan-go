package services

import (
	"io"
	"pan/app/constant"
	"pan/app/node"
	"pan/extfs/models"

	"google.golang.org/protobuf/proto"
)

type RemoteNodeItemInternalService interface {
	TraverseAllWithNodeID(traverseFn func(node *models.RemoteNodeItem) error, nodeId node.NodeID) error
}

type RemoteNodeItemService struct {
	NodeModule node.NodeModule
}

var RequestAllRemoteItems = []byte("extfs/select_all_remote_items")

func (s *RemoteNodeItemService) TraverseAllWithNodeID(traverseFn func(node *models.RemoteNodeItem) error, nodeId node.NodeID) error {
	request := node.NewRequest(RequestAllRemoteItems, nil)

	response, err := s.NodeModule.Do(nodeId, request)
	if err != nil {
		return err
	}
	if response.Code() != constant.CodeOK {
		return constant.ErrInternalError
	}
	data, err := io.ReadAll(response.Body())
	if err != nil {
		return err
	}

	var remoteNodeItemList models.RemoteNodeItemList
	err = proto.Unmarshal(data, &remoteNodeItemList)
	if err != nil {
		return err
	}

	for _, item := range remoteNodeItemList.Items {
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}
