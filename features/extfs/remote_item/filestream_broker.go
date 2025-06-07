package remoteitem

import (
	"context"
	"io"
	"pan/lib/feature"
	"pan/lib/peer"
)

var SelectRemoteFileStream = []byte("select_remote_filestream")

type RemoteFileStreamBroker struct {
	feature.BaseBroker
}

func (broker *RemoteFileStreamBroker) Select(peerId peer.PeerID, condition *RemoteFileStreamSelectCondition) (io.ReadCloser, error) {
	return broker.Request(context.Background(), peerId, SelectRemoteFileStream, condition)
}
