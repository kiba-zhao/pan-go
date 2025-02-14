package remoteitem

import (
	"bytes"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteFileInfoTopic struct {
	RemoteFileInfoService *RemoteFileInfoService
}

func (topic *RemoteFileInfoTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SearchRemoteFileInfos, topic.Search)
	router.Handle(SelectRemoteFileInfo, topic.Select)
	return nil
}

func (topic *RemoteFileInfoTopic) Search(ctx *peer.Context, next peer.Next) error {

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteFileInfoRecordSearchCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	fileItemList, err := topic.RemoteFileInfoService.SearchForTopic(&condition)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	buffer, err := proto.Marshal(fileItemList)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	ctx.Respond(bytes.NewReader(buffer))
	return err
}

func (topic *RemoteFileInfoTopic) Select(ctx *peer.Context, next peer.Next) error {

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteFileInfoRecordSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	record, err := topic.RemoteFileInfoService.SelectForTopic(&condition)
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
