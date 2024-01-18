package peer

import (
	"pan/modules/extfs/models"
	"pan/peer"
)

type API interface {
	GetPeerInfo(peerId peer.PeerId) (models.PeerInfo, error)
	// GetFileInfos(peerId peer.PeerId, hash []byte) (models.FileInfo, error)
}

type apiImpl struct {
	peer peer.Peer
}

// // NewPeerAPI ...
// func NewPeerAPI(p peer.Peer) PeerAPI {
// 	api := new(peerAPI)
// 	api.peer = p
// 	return api
// }
