package peer

import (
	"io"
	"pan/modules/extfs/models"
	"pan/peer"

	"google.golang.org/protobuf/proto"
)

type API interface {
	GetRemoteFilesState(peerId peer.PeerId) (models.RemoteStateInfo, error)
}

type apiImpl struct {
	Peer peer.Peer
}

func NewAPI(peer peer.Peer) API {
	api := new(apiImpl)
	api.Peer = peer
	return api
}

func (a *apiImpl) GetRemoteFilesState(peerId peer.PeerId) (models.RemoteStateInfo, error) {

	node, err := a.Peer.Open(peerId)
	if err != nil {
		return models.RemoteStateInfo{}, err
	}

	res, err := a.Peer.Request(node, nil, []byte("GetRemoteFilesState"))
	if err != nil {
		return models.RemoteStateInfo{}, err
	}

	body, err := io.ReadAll(res.Body())
	if err != nil {
		return models.RemoteStateInfo{}, err
	}

	var info models.RemoteStateInfo
	err = proto.Unmarshal(body, &info)
	return info, err
}
