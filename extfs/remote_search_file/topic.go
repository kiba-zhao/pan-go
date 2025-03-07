// Define topic for remote search file
//
// It implements then api for peer service
package remotesearchfile

import (
	"bytes"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteSearchFileTopic struct {
	RemoteSearchFileService *RemoteSearchFileService
}

// SetupToPeer sets up the topic to peer service
//
// It sets up the following endpoint:
//
// - `SearchRemoteSearchFiles`: Search all files with the given condition
func (t *RemoteSearchFileTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SearchRemoteSearchFiles, t.Search)
	return nil
}

// Search is called when a peer request with `SearchRemoteSearchFiles` topic is received
//
// It reads the request body, unmarshals it into a RemoteSearchFileRecordSearchCondition, and
// calls the SearchForTopic method of the RemoteSearchFileService to search for remote search files
// with the given condition.
//
// If the request is invalid, it throws an error with CodeBadRequest.
//
// If the search request fails, it returns the error.
//
// Otherwise, it marshals the result into a RemoteSearchFileRecordList and responds with the marshaled
// data.
func (t *RemoteSearchFileTopic) Search(ctx *peer.Context, next peer.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	var condition RemoteSearchFileRecordSearchCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	list, err := t.RemoteSearchFileService.SearchForTopic(&condition)
	if err != nil {
		return err
	}

	buffer, err := proto.Marshal(list)
	if err == nil {
		ctx.Respond(bytes.NewReader(buffer))
	}

	return err
}
