package remoteitem

import (
	"errors"
	"io"
	nodeitem "pan/features/extfs/node_item"
	"pan/lib/peer"

	"google.golang.org/protobuf/proto"
)

type RemoteFileStreamTopic struct {
	RemoteFileStreamService *RemoteFileStreamService
}

func (topic *RemoteFileStreamTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SelectRemoteFileStream, topic.Select)
	return nil
}

func (topic *RemoteFileStreamTopic) Select(ctx peer.PeerContext, next peer.PeerNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	var condition RemoteFileStreamSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	res, err := topic.RemoteFileStreamService.SelectForTopic(&condition)
	if errors.Is(err, nodeitem.ErrNodeFilePathWithoutFolder) {
		ctx.ThrowError(peer.CodeForbidden, err)
		return nil
	}

	if err == nil {
		ctx.Respond(res)
	}

	return err
}
