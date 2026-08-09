package user

import (
	"context"
	"pan/pkg/proto"
	"pan/pkg/ptp"
)

type UserDataTopic struct {
	UserDataService *UserDataService
}

var _ = (ptp.PeerTopic)((*UserDataTopic)(nil))

func (topic *UserDataTopic) SetupToPeer(router ptp.PeerServletRouter) error {
	router.Handle(PullUserData, topic.Pull)
	router.Handle(PushUserData, topic.Push)
	return nil
}

func (topic *UserDataTopic) Pull(ctx ptp.PeerServletContext, next ptp.PeerServletNext) error {
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
	user, err := topic.UserDataService.PullForTopic(context.Background(), peerId.(ptp.PeerID), meta)
	if err != nil {
		return err
	}

	remoteUser := parseRemoteUser(user)
	reader, err := proto.MarshalWithReader(remoteUser)
	if err == nil {
		ctx.Respond(reader)
	}

	return err
}

func (topic *UserDataTopic) Push(ctx ptp.PeerServletContext, next ptp.PeerServletNext) error {
	peerId, ok := ctx.Session(ptp.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(ptp.CodeBadRequest, nil)
		return nil
	}

	signature, ok := ctx.RequestHeader(PeerSignatureHeaderName)
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
	err = topic.UserDataService.PushForTopic(context.Background(), peerId.(ptp.PeerID), meta, signature)
	return err
}
