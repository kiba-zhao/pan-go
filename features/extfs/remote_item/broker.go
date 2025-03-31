// Define Broker for remote item
//
// It is responsible for searching remote items by peer request
package remoteitem

import (
	"context"
	"pan/lib/peer"
	"pan/lib/sample"
)

var SelectAllRemoteItems = []byte("select_all_remote_items")
var SelectRemoteItem = []byte("select_remote_item")

type RemoteItemBroker struct {
	SamplePeer sample.SamplePeer
}

// SelectAll sends a request to select all remote items with the given peer ID.
//
// Parameters:
//   - peerId: The ID of the peer to which the request is sent.
//
// Returns:
//   - A pointer to a RemoteItemRecordList containing the list of all remote items.
//   - An error if the request fails.
func (broker *RemoteItemBroker) SelectAll(peerId peer.PeerID) (*RemoteItemRecordList, error) {

	var remoteItemRecordList RemoteItemRecordList
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SelectAllRemoteItems, &remoteItemRecordList, nil)

	return &remoteItemRecordList, err
}

// Select sends a request to select a remote item with the given peer ID and condition.
//
// Parameters:
//   - peerId: The ID of the peer to which the request is sent.
//   - condition: The condition used to select the remote item.
//
// Returns:
//   - A pointer to a RemoteItemRecord containing the selected remote item.
//   - An error if the request fails.
func (broker *RemoteItemBroker) Select(peerId peer.PeerID, condition *RemoteItemRecordSelectCondition) (*RemoteItemRecord, error) {

	var record RemoteItemRecord
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SelectRemoteItem, &record, condition)

	return &record, err
}
