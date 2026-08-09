package user

import (
	"context"
	"pan/pkg/proto"
	"pan/pkg/ptp"
)

type UserSecretTopic struct {
	UserSecretService *UserSecretService
}

var _ = (ptp.PeerTopic)((*UserSecretTopic)(nil))

func (topic *UserSecretTopic) SetupToPeer(router ptp.PeerServletRouter) error {
	router.Handle(PullUserSecret, topic.Pull)
	return nil
}

func (topic *UserSecretTopic) Pull(ctx ptp.PeerServletContext, next ptp.PeerServletNext) error {
	peerId, ok := ctx.Session(ptp.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(ptp.CodeBadRequest, nil)
		return nil
	}

	var userMeta RemoteUserMeta
	err := proto.UnmarshalWithReader(ctx.Request(), &userMeta)
	if err != nil {
		ctx.ThrowError(ptp.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)
	secret, err := topic.UserSecretService.PullForTopic(context.Background(), peerId.(ptp.PeerID), meta)
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
