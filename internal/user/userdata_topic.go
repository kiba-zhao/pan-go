package user

import (
	"bytes"
	"io"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

type UserDataTopic struct {
	UserDataService *UserDataService
}

func (topic *UserDataTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(PullUserData, topic.Pull)
	router.Handle(PushUserData, topic.Push)
	return nil
}

func (topic *UserDataTopic) Pull(ctx peer.PeerContext, next peer.PeerNext) error {
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
	user, err := topic.UserDataService.PullForTopic(peerId.(peer.PeerID), meta)
	if err != nil {
		return err
	}

	remoteUser := parseRemoteUser(user)
	buffer, err := proto.Marshal(remoteUser)
	if err == nil {
		ctx.Respond(bytes.NewReader(buffer))
	}

	return err
}

func (topic *UserDataTopic) Push(ctx peer.PeerContext, next peer.PeerNext) error {
	peerId, ok := ctx.Session(peer.ContextPeerID)
	if !ok {
		ctx.ThrowError(peer.CodeBadRequest, nil)
		return nil
	}

	signature, ok := ctx.RequestHeader(PeerSignatureHeaderName)
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

	err = topic.UserDataService.PushForTopic(peerId.(peer.PeerID), meta, signature)
	return err
}
