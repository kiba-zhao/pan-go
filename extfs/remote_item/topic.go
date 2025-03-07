// Define topic for remote item
//
// It implements then api for peer service
package remoteitem

import (
	"bytes"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteItemTopic struct {
	RemoteItemService *RemoteItemService
}

// SetupToPeer sets up the topic to peer service
//
// It sets up the following endpoints:
//
// - `SelectAllRemoteItems`: Select all remote items with the given peer ID
// - `SelectRemoteItem`: Select a remote item by its ID and peer ID
func (topic *RemoteItemTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SelectAllRemoteItems, topic.SelectAll)
	router.Handle(SelectRemoteItem, topic.Select)
	return nil
}

// SelectAll is called when a peer request with `SelectAllRemoteItems` topic is received
//
// It calls the SelectAllForTopic method of the RemoteItemService to select all remote items
// with the given peer ID.
//
// If the request is invalid, it throws an error with CodeBadRequest.
//
// If the select request fails, it returns the error.
//
// Otherwise, it marshals the result into a RemoteItemRecordList and responds with the marshaled
// data.
func (topic *RemoteItemTopic) SelectAll(ctx *peer.Context, next peer.Next) error {

	recordList, err := topic.RemoteItemService.SelectAllForTopic()
	if err != nil {
		return err
	}

	buffer, err := proto.Marshal(&recordList)
	if err == nil {
		ctx.Respond(bytes.NewReader(buffer))
	}

	return err
}

func (topic *RemoteItemTopic) Select(ctx *peer.Context, next peer.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	var condition RemoteItemRecordSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	record, err := topic.RemoteItemService.SelectForTopic(&condition)
	if err != nil {
		return err
	}

	buffer, err := proto.Marshal(record)
	if err == nil {
		ctx.Respond(bytes.NewReader(buffer))
	}

	return err
}
