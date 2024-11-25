package remotefile

import (
	"pan/app/peer"
	"pan/app/sample"
)

var SearchRemoteFiles = []byte("search_remote_file_items")
var SelectRemoteFile = []byte("select_remote_file_item")

type RemoteFileBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteFileBroker) Search(peerId peer.PeerID, condition *RemoteFileRecordSearchCondition) (*RemoteFileRecordList, error) {
	var remoteFileRecordList RemoteFileRecordList
	err := broker.SamplePeer.RequestWithProto(peerId, SearchRemoteFiles, &remoteFileRecordList, condition)

	return &remoteFileRecordList, err
}

func (broker *RemoteFileBroker) Select(peerId peer.PeerID, condition *RemoteFileRecordSelectCondition) (*RemoteFileRecord, error) {
	var remoteFileRecord RemoteFileRecord
	err := broker.SamplePeer.RequestWithProto(peerId, SelectRemoteFile, &remoteFileRecord, condition)

	return &remoteFileRecord, err
}
