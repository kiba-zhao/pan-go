package remotefile

import (
	"bytes"
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteFileTopic struct {
	RemoteFileService *RemoteFileService
}

func (topic *RemoteFileTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SearchRemoteFiles, topic.Search)
	router.Handle(SelectRemoteFile, topic.Select)
	return nil
}

func (topic *RemoteFileTopic) Search(ctx *peer.Context, next peer.Next) error {

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteFileRecordSearchCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	fileItemList, err := topic.RemoteFileService.SearchForTopic(&condition)
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

func (topic *RemoteFileTopic) Select(ctx *peer.Context, next peer.Next) error {

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteFileRecordSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	record, err := topic.RemoteFileService.SelectForTopic(&condition)
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
