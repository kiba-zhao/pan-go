package peer_test

import (
	"bytes"
	"pan/modules/extfs/models"
	"pan/modules/extfs/peer"
	corePeer "pan/peer"
	"testing"

	mockedPeer "pan/mocks/pan/peer"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// create unit test for peer API
func TestAPI(t *testing.T) {

	// setup function has one argument of type corePeer.Peer and return API interface.
	setup := func(p corePeer.Peer) peer.API {
		api := peer.NewAPI(p)
		return api
	}

	t.Run("GetPeerInfo", func(t *testing.T) {

		p := new(mockedPeer.MockPeer)
		api := setup(p)

		peerId := uuid.New()
		node := new(mockedPeer.MockNode)
		info := new(models.PeerInfo)
		info.PeerId = peerId[:]
		info.Hash = []byte("hash")
		info.Time = 123
		bodyBytes, err := proto.Marshal(info)
		if err != nil {
			t.Fatal(err)
		}
		reader := bytes.NewReader(bodyBytes)
		res := corePeer.NewResponse(200, reader)

		p.On("Open", peerId).Return(node, nil)
		p.On("Request", node, nil, []byte("GetPeerInfo")).Return(res, nil)

		resInfo, err := api.GetPeerInfo(peerId)

		assert.Nil(t, err)
		assert.Equal(t, info.PeerId, resInfo.PeerId)
		assert.Equal(t, info.Hash, resInfo.Hash)
		assert.Equal(t, info.Time, resInfo.Time)

		p.AssertExpectations(t)
		node.AssertExpectations(t)
	})
}
