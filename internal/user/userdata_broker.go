package user

import (
	"context"
	"pan/internal/feature"
	"pan/internal/peer"
)

type UserDataBroker struct {
	PeerBroker *feature.PeerBroker
}

func (broker *UserDataBroker) Pull(peerId peer.PeerID, meta *RemoteUserMeta) (*RemoteUser, error) {
	ctx := context.Background()
	var remoteUser RemoteUser

	err := broker.PeerBroker.RequestWithProto(ctx, peerId, PullUserData, &remoteUser, meta)

	return &remoteUser, err
}

func (broker *UserDataBroker) Push(peerId peer.PeerID, meta *RemoteUserMeta, signature []byte) error {
	ctx := context.Background()

	signatureHeader := peer.PeerHeaderItem{Key: PeerSignatureHeaderName, Value: signature}
	res, err := broker.PeerBroker.Request(ctx, peerId, PushUserData, meta, signatureHeader)
	if err == nil {
		res.Close()
	}
	return err
}
