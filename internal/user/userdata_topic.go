package user

import (
	"pan/pkg/net"
	"pan/pkg/proto"
)

type UserDataTopic struct {
	UserDataService *UserDataService
}

var _ = (net.PeerTopic)((*UserDataTopic)(nil))

func (topic *UserDataTopic) SetupToPeer(router net.PeerServletRouter) error {
	router.Handle(PullUserData, topic.Pull)
	router.Handle(PushUserData, topic.Push)
	return nil
}

func (topic *UserDataTopic) Pull(ctx net.PeerServletContext, next net.PeerServletNext) error {
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
	user, err := topic.UserDataService.PullForTopic(peerId.(net.PeerID), meta)
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

func (topic *UserDataTopic) Push(ctx net.PeerServletContext, next net.PeerServletNext) error {
	peerId, ok := ctx.Session(net.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(net.CodeBadRequest, nil)
		return nil
	}

	signature, ok := ctx.RequestHeader(PeerSignatureHeaderName)
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
	err = topic.UserDataService.PushForTopic(peerId.(net.PeerID), meta, signature)
	return err
}
