// Define Broker for remote search file
//
// It is responsible for searching remote search files by peer request.
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

// Search sends a request to search for remote search files based on the given peer ID and search condition.
// It utilizes the SamplePeer to perform the request and returns a list of remote search file records or an error if the request fails.
//
// Parameters:
//   - peerId: The ID of the peer to which the search request is sent.
//   - condition: The search condition used to filter the remote search files.
//
// Returns:
//   - A pointer to a RemoteSearchFileRecordList containing the search results.
//   - An error if the search request fails.

func (broker *RemoteSearchFileBroker) Search(peerId peer.PeerID, condition *RemoteSearchFileRecordSearchCondition) (*RemoteSearchFileRecordList, error) {
	var list RemoteSearchFileRecordList
	err := broker.SamplePeer.RequestWithProto(context.Background(), peerId, SearchRemoteSearchFiles, &list, condition)
	return &list, err
}
