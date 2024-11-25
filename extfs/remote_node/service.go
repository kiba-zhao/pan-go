package remotenode

import (
	"encoding/base64"
	"errors"
	"slices"

	appnode "pan/app/app_node"
	"pan/app/peer"
)

var ErrRemoteNodeUnavailable = errors.New("remotenode.RemoteNodeService Error: Unavailable")
var ErrRemoteNodeInvalidPeer = errors.New("remotenode.RemoteNodeService Error: Invalid Peer")

type RemoteNodeService struct {
	PeerManager            peer.PeerManager
	AppNodeExternalService appnode.AppNodeExternalService
}

func (s *RemoteNodeService) SelectAll() (int64, []RemoteNode, error) {
	mgr := s.PeerManager
	if mgr == nil {
		return 0, nil, ErrRemoteNodeUnavailable
	}

	peerIds := make([]string, 0)
	err := mgr.TraversePeerID(func(peerId peer.PeerID) error {
		peerIdBase64 := base64.StdEncoding.EncodeToString(peerId)
		idx, _ := slices.BinarySearch(peerIds, peerIdBase64)
		peerIds = slices.Insert(peerIds, idx, peerIdBase64)
		return nil
	})

	if err != nil || len(peerIds) <= 0 {
		return 0, nil, err
	}

	remotes := make([]RemoteNode, 0)
	err = s.AppNodeExternalService.TraverseWithPeerIDs(func(model appnode.AppNode) error {
		idx, ok := slices.BinarySearch(peerIds, model.PeerID)
		if !ok {
			return ErrRemoteNodeInvalidPeer
		}
		peerIds = slices.Delete(peerIds, idx, idx+1)

		var remote RemoteNode
		remote.ID = model.ID
		remote.Name = model.Name
		remote.PeerID = model.PeerID
		remote.Available = true
		remote.TagQuantity = 0
		remote.PendingTagQuantity = 0
		remote.CreatedAt = model.CreatedAt
		remote.UpdatedAt = model.UpdatedAt

		// TODO: set tag quantity

		remotes = append(remotes, remote)

		return nil
	}, peerIds)

	return 0, remotes, err
}

func (s *RemoteNodeService) SelectByName(name string) (RemoteNode, error) {
	model, err := s.AppNodeExternalService.SelectByName(name)
	var remote RemoteNode
	if err != nil {
		return remote, err
	}
	remote.ID = model.ID
	remote.Name = model.Name
	remote.PeerID = model.PeerID
	remote.Available = !model.Blocked && model.Online
	remote.TagQuantity = 0
	remote.PendingTagQuantity = 0
	remote.CreatedAt = model.CreatedAt
	remote.UpdatedAt = model.UpdatedAt

	// TODO: set tag quantity
	return remote, err
}
