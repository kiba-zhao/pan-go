package user

import (
	"context"
	"pan/pkg/proto"
	"pan/pkg/ptp"
)

type UserSecretBroker struct {
	PeerBroker *ptp.PeerBroker
}

func (broker *UserSecretBroker) Pull(ctx context.Context, peerId ptp.PeerID, meta *RemoteUserMeta) (*RemoteUserSecret, error) {
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
