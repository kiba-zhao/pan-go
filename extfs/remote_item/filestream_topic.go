package remoteitem

import (
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteFileStreamTopic struct {
	RemoteFileStreamService *RemoteFileStreamService
}

func (topic *RemoteFileStreamTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SelectRemoteFileStream, topic.Select)
	return nil
}

func (topic *RemoteFileStreamTopic) Select(ctx *peer.Context, next peer.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteFileStreamSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	res, err := topic.RemoteFileStreamService.SelectForTopic(&condition)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	ctx.Respond(res)
	return err
}
