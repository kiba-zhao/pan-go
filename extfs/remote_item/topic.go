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

func (topic *RemoteItemTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SelectAllRemoteItems, topic.SelectAll)
	router.Handle(SelectRemoteItem, topic.Select)
	return nil
}

func (topic *RemoteItemTopic) SelectAll(ctx *peer.Context, next peer.Next) error {

	recordList, err := topic.RemoteItemService.SelectAllForTopic()
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	buffer, err := proto.Marshal(&recordList)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	ctx.Respond(bytes.NewReader(buffer))
	return err
}

func (topic *RemoteItemTopic) Select(ctx *peer.Context, next peer.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteItemRecordSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	record, err := topic.RemoteItemService.SelectForTopic(&condition)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	buffer, err := proto.Marshal(record)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	ctx.Respond(bytes.NewReader(buffer))
	return err
}
