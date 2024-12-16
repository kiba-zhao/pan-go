package remotenode

import (
	appnode "pan/app/app_node"
)

type RemoteNodeService struct {
	AppNodeExternalService appnode.AppNodeExternalService
}

func (s *RemoteNodeService) SelectAll() (int64, []RemoteNode, error) {

	remotes := make([]RemoteNode, 0)
	err := s.AppNodeExternalService.TraverseAll(func(model appnode.AppNode) error {

		remote := parseRemoteNode(model)
		if remote.Available {
			// TODO: set tag quantity
			remotes = append(remotes, remote)
		}

		return nil
	})

	return 0, remotes, err
}

func (s *RemoteNodeService) SelectByName(name string) (RemoteNode, error) {
	model, err := s.AppNodeExternalService.SelectByName(name)
	if err != nil {
		return RemoteNode{}, err
	}

	remote := parseRemoteNode(model)
	// TODO: set tag quantity
	return remote, err
}

func parseRemoteNode(model appnode.AppNode) RemoteNode {
	var remote RemoteNode
	remote.ID = model.ID
	remote.Name = model.Name
	remote.PeerID = model.PeerID
	remote.Available = !model.Blocked && model.Online
	remote.TagQuantity = 0
	remote.PendingTagQuantity = 0
	remote.CreatedAt = model.CreatedAt
	remote.UpdatedAt = model.UpdatedAt
	return remote
}
