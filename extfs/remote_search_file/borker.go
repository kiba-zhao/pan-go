package remotesearchfile

import (
	"context"
	"pan/app/peer"
	"pan/app/sample"
)

var SearchRemoteSearchFiles = []byte("search_remote_search_files")

type RemoteSearchFileBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteSearchFileBroker) Search(peerId peer.PeerID, condition *RemoteSearchFileRecordSearchCondition) (*RemoteSearchFileRecordList, error) {
	var list RemoteSearchFileRecordList
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SearchRemoteSearchFiles, &list, condition)
	return &list, err
}
