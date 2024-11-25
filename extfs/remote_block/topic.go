package remoteblock

import (
	"io"
	"pan/app/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteBlockTopic struct {
	RemoteBlockService *RemoteBlockService
}

func (topic *RemoteBlockTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SelectRemoteBlock, topic.Select)
	return nil
}

func (topic *RemoteBlockTopic) Select(ctx *peer.Context, next peer.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	var condition RemoteBlockSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return err
	}

	res, err := topic.RemoteBlockService.SelectForTopic(&condition)
	if err != nil {
		ctx.ThrowError(peer.CodeInternalError, err)
		return err
	}

	ctx.Respond(res)
	return err
}
