package user

import (
	"context"
	"pan/pkg/net"
	"pan/pkg/proto"
)

type UserSecretBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserSecretBroker) Pull(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta) (*RemoteUserSecret, error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(PullUserSecret, reader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if res != nil {
		defer res.Close()
	}
	if err != nil {
		return nil, err
	}

	var secret RemoteUserSecret
	err = proto.UnmarshalWithReader(res, &secret)

	return &secret, err
}
