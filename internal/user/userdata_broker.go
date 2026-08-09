package user

import (
	"context"
	"pan/pkg/proto"
	"pan/pkg/ptp"
)

type UserDataBroker struct {
	PeerBroker *ptp.PeerBroker
}

func (broker *UserDataBroker) Pull(ctx context.Context, peerId ptp.PeerID, meta *RemoteUserMeta) (*RemoteUser, error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(PullUserData, reader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if res != nil {
		defer res.Close()
	}
	if err != nil {
		return nil, err
	}

	var remoteUser RemoteUser
	err = proto.UnmarshalWithReader(res, &remoteUser)

	return &remoteUser, err
}

func (broker *UserDataBroker) Push(ctx context.Context, peerId ptp.PeerID, meta *RemoteUserMeta, signature []byte) error {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return err
	}

	signatureHeader := ptp.HeaderItem{Key: PeerSignatureHeaderName, Value: signature}
	req := broker.PeerBroker.NewRequest(PushUserData, reader, signatureHeader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if res != nil {
		defer res.Close()
	}
	return err
}
