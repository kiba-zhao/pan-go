package remoteitem

import (
	"context"
	"io"
	"pan/app/peer"
	"pan/app/sample"
)

var SelectRemoteFileStream = []byte("select_remote_filestream")

type RemoteFileStreamBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteFileStreamBroker) Select(peerId peer.PeerID, condition *RemoteFileStreamSelectCondition) (io.ReadCloser, error) {
	return broker.SamplePeer.Request(context.Background(), peerId, SelectRemoteFileStream, condition)
}
