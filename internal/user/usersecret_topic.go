package user

import (
	"bytes"
	"io"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

type UserSecretTopic struct {
	UserSecretService *UserSecretService
}

func (topic *UserSecretTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(PullUserSecret, topic.Pull)
	return nil
}

func (topic *UserSecretTopic) Pull(ctx peer.PeerContext, next peer.PeerNext) error {
	peerId, ok := ctx.Session(peer.ContextPeerID)
	if !ok {
		ctx.ThrowError(peer.CodeBadRequest, nil)
		return nil
	}

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	var userMeta RemoteUserMeta
	err = proto.Unmarshal(body, &userMeta)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)
	secret, err := topic.UserSecretService.PullForTopic(peerId.(peer.PeerID), meta)
	if err != nil {
		return err
	}

	remoteSecret := parseRemoteUserSecret(secret)
	buffer, err := proto.Marshal(remoteSecret)
	if err == nil {
		ctx.Respond(bytes.NewReader(buffer))
	}
	return err
}
