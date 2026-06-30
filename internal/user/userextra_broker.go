package user

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"iter"
	"pan/pkg/net"
	"pan/pkg/proto"

	protobuf "google.golang.org/protobuf/proto"
)

var ErrUserExtraScanInvalidWithRemoteUser = errors.New("user.UserExtraBroker ScanWithUserMeta Error: Invalid")

type UserExtraBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserExtraBroker) ScanWithUserMeta(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserExtra, error], error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(ScanUserExtraWithUserMeta, reader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if err != nil {
		if res != nil {
			defer res.Close()
		}
		return nil, err
	}

	return func(yield func(*RemoteUserExtra, error) bool) {
		defer res.Close()

		sizeBytes := make([]byte, 2)
		var itemBytes []byte
		for {
			n, err := res.Read(sizeBytes)
			if err == nil && n != 2 {
				err = ErrUserExtraScanInvalidWithRemoteUser
			}
			if err == nil {
				size := binary.BigEndian.Uint16(sizeBytes)
				itemBytes = make([]byte, size)
				if size > 0 {
					n, err = res.Read(itemBytes)
					if err == nil && n != int(size) {
						err = ErrUserExtraScanInvalidWithRemoteUser
					}
				}
			}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					yield(nil, err)
				}
				break
			}

			var remoteExtra RemoteUserExtra
			err = protobuf.Unmarshal(itemBytes, &remoteExtra)
			if !yield(&remoteExtra, err) {
				break
			}
		}
	}, nil
}

func parseUserExtraSeq(remoteExtraSeq iter.Seq2[*RemoteUserExtra, error]) iter.Seq2[UserExtra, error] {
	return func(yield func(UserExtra, error) bool) {
		for remoteExtra, err := range remoteExtraSeq {
			var extra UserExtra
			if err == nil {
				extra = parseUserExtra(remoteExtra)
			}

			if !yield(extra, err) {
				break
			}
		}
	}
}
