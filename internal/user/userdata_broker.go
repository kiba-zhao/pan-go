package user

import (
	"context"
	"pan/pkg/net"
	"pan/pkg/proto"
)

type UserDataBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserDataBroker) Pull(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta) (*RemoteUser, error) {
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

func (broker *UserDataBroker) Push(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta, signature []byte) error {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return err
	}

	signatureHeader := net.HeaderItem{Key: PeerSignatureHeaderName, Value: signature}
	req := broker.PeerBroker.NewRequest(PushUserData, reader, signatureHeader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if res != nil {
		defer res.Close()
	}
	return err
}
