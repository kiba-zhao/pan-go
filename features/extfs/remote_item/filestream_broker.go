package remoteitem

import (
	"context"
	"io"
	"pan/lib/peer"
	"pan/lib/sample"
)

var SelectRemoteFileStream = []byte("select_remote_filestream")

type RemoteFileStreamBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteFileStreamBroker) Select(peerId peer.PeerID, condition *RemoteFileStreamSelectCondition) (io.ReadCloser, error) {
	return broker.SamplePeer.Request(context.Background(), peerId, SelectRemoteFileStream, condition)
}
