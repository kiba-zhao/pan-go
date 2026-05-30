package user

import (
	"pan/internal/net"
	"pan/internal/proto"
)

type UserSecretTopic struct {
	UserSecretService *UserSecretService
}

var _ = (net.PeerTopic)((*UserSecretTopic)(nil))

func (topic *UserSecretTopic) SetupToPeer(router net.PeerServletRouter) error {
	router.Handle(PullUserSecret, topic.Pull)
	return nil
}

func (topic *UserSecretTopic) Pull(ctx net.PeerServletContext, next net.PeerServletNext) error {
	peerId, ok := ctx.Session(net.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(net.CodeBadRequest, nil)
		return nil
	}

	var userMeta RemoteUserMeta
	err := proto.UnmarshalWithReader(ctx.Request(), &userMeta)
	if err != nil {
		ctx.ThrowError(net.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)
	secret, err := topic.UserSecretService.PullForTopic(peerId.(net.PeerID), meta)
	if err != nil {
		return err
	}

	remoteSecret := parseRemoteUserSecret(secret)
	reader, err := proto.MarshalWithReader(remoteSecret)
	if err == nil {
		ctx.Respond(reader)
	}
	return err
}
