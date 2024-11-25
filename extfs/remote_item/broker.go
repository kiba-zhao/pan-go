package remoteitem

import (
	"pan/app/peer"
	"pan/app/sample"
)

var SelectAllRemoteItems = []byte("select_all_remote_items")
var SelectRemoteItem = []byte("select_remote_item")

type RemoteItemBroker struct {
	SamplePeer sample.SamplePeer
}

func (broker *RemoteItemBroker) SelectAll(peerId peer.PeerID) (*RemoteItemRecordList, error) {

	var remoteItemRecordList RemoteItemRecordList
	err := broker.SamplePeer.RequestWithProto(peerId, SelectAllRemoteItems, &remoteItemRecordList, nil)

	return &remoteItemRecordList, err
}

func (broker *RemoteItemBroker) Select(peerId peer.PeerID, condition *RemoteItemRecordSelectCondition) (*RemoteItemRecord, error) {

	var record RemoteItemRecord
	err := broker.SamplePeer.RequestWithProto(peerId, SelectRemoteItem, &record, condition)

	return &record, err
}
