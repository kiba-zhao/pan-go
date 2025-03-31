package remoteitem

import (
	"context"
	"pan/lib/peer"
	"pan/lib/sample"
)

var SearchRemoteFileInfos = []byte("search_remote_fileinfos")
var SelectRemoteFileInfo = []byte("select_remote_fileinfo")

type RemoteFileInfoBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteFileInfoBroker) Search(peerId peer.PeerID, condition *RemoteFileInfoRecordSearchCondition) (*RemoteFileInfoRecordList, error) {
	var list RemoteFileInfoRecordList
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SearchRemoteFileInfos, &list, condition)

	return &list, err
}

func (broker *RemoteFileInfoBroker) Select(peerId peer.PeerID, condition *RemoteFileInfoRecordSelectCondition) (*RemoteFileInfoRecord, error) {
	var record RemoteFileInfoRecord
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SelectRemoteFileInfo, &record, condition)

	return &record, err
}
