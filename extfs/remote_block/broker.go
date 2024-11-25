package remoteblock

import (
	"io"
	"pan/app/peer"
	"pan/app/sample"
)

var SelectRemoteBlock = []byte("select_remote_block")

type RemoteBlockBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteBlockBroker) Select(peerId peer.PeerID, condition *RemoteBlockSelectCondition) (io.ReadCloser, error) {
	return broker.SamplePeer.Request(peerId, SelectRemoteBlock, condition)
}
