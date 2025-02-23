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

func (t *RemoteSearchFileTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(SearchRemoteSearchFiles, t.Search)
	return nil
}

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
