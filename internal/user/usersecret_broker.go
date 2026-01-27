package user

import (
	"context"
	"pan/internal/feature"
	"pan/internal/peer"
)

type UserSecretBroker struct {
	PeerBroker *feature.PeerBroker
}

func (broker *UserSecretBroker) Pull(peerId peer.PeerID, meta *RemoteUserMeta) (*RemoteUserSecret, error) {
	ctx := context.Background()
	var secret RemoteUserSecret

	err := broker.PeerBroker.RequestWithProto(ctx, peerId, PullUserSecret, &secret, meta)
	return &secret, err
}
