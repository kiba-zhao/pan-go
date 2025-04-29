package node_test

import (
	"iter"
	appnode "pan/features/app/node"
	"pan/lib/peer"
	mocked "pan/mocks/pan/features/app/node"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkAddrGuide(t *testing.T) {

	setup := func() *appnode.NetworkAddrGuide {
		var guide appnode.NetworkAddrGuide
		guide.NetworkAddrService = &appnode.NetworkAddrService{}
		return &guide
	}

	t.Run("LookupExplorerAddr", func(t *testing.T) {
		guide := setup()

		// mock NetworkAddrRepository
		networkAddrRepo := &mocked.MockNetworkAddrRepository{}
		defer networkAddrRepo.AssertExpectations(t)
		guide.NetworkAddrService.NetworkAddrRepo = networkAddrRepo

		peerId := peer.PeerID{1, 2, 3}
		peerIdStr := peer.EncodePeerID(peerId)

		var entity appnode.NetworkAddr
		entity.AppNodeID = 123
		entity.Address = "address"
		entity.ID = 1

		fakeSeq := fakeNetworkAddrSeq(entity)

		networkAddrRepo.On("SeqByPeerId", peerIdStr).Once().Return(fakeSeq, nil)

		seq, err := guide.LookupExplorerAddr(peerId)
		if err != nil {
			t.Fatal(err)
		}

		count := 0
		for addr := range seq {
			assert.Equal(t, entity.Address, addr)
			count++
		}

		assert.Equal(t, 1, count)
	})

}

func fakeNetworkAddrSeq(addrs ...appnode.NetworkAddr) iter.Seq[appnode.NetworkAddr] {
	return func(yield func(appnode.NetworkAddr) bool) {
		for _, addr := range addrs {
			if !yield(addr) {
				break
			}
		}
	}
}
