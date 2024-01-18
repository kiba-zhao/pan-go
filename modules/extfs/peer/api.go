package peer

import (
	"io"
	"pan/modules/extfs/models"
	"pan/peer"

	"google.golang.org/protobuf/proto"
)

type API interface {
	GetPeerInfo(peerId peer.PeerId) (models.PeerInfo, error)
	// GetFileInfos(peerId peer.PeerId, hash []byte) (models.FileInfo, error)
}

type apiImpl struct {
	Peer peer.Peer
}

func NewAPI(peer peer.Peer) API {
	api := new(apiImpl)
	api.Peer = peer
	return api
}

func (a *apiImpl) GetPeerInfo(peerId peer.PeerId) (models.PeerInfo, error) {

	node, err := a.Peer.Open(peerId)
	if err != nil {
		return models.PeerInfo{}, err
	}

	res, err := a.Peer.Request(node, nil, []byte("GetPeerInfo"))
	if err != nil {
		return models.PeerInfo{}, err
	}

	body, err := io.ReadAll(res.Body())
	if err != nil {
		return models.PeerInfo{}, err
	}

	var info models.PeerInfo
	err = proto.Unmarshal(body, &info)
	return info, err
}
