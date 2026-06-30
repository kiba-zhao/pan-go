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

var ErrUserConsensusScanInvalidWithRemoteUser = errors.New("user.UserConsensusBroker ScanWithUserMeta Error: Invalid")

type UserConsensusBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserConsensusBroker) ScanWithUserMeta(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserConsensus, error], error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(ScanUserConsensusWithUserMeta, reader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
	if err != nil {
		if res != nil {
			defer res.Close()
		}
		return nil, err
	}

	return func(yield func(*RemoteUserConsensus, error) bool) {
		defer res.Close()

		sizeBytes := make([]byte, 2)
		var itemBytes []byte
		for {
			n, err := res.Read(sizeBytes)
			if err == nil && n != 2 {
				err = ErrUserConsensusScanInvalidWithRemoteUser
			}
			if err == nil {
				size := binary.BigEndian.Uint16(sizeBytes)
				itemBytes = make([]byte, size)
				if size > 0 {
					n, err = res.Read(itemBytes)
					if err == nil && n != int(size) {
						err = ErrUserConsensusScanInvalidWithRemoteUser
					}
				}
			}
			if err != nil {
				if !errors.Is(err, io.EOF) {
					yield(nil, err)
				}
				break
			}

			var remoteConsensus RemoteUserConsensus
			err = protobuf.Unmarshal(itemBytes, &remoteConsensus)
			if !yield(&remoteConsensus, err) {
				break
			}
		}
	}, nil
}

func parseUserConsensusSeq(seq iter.Seq2[*RemoteUserConsensus, error]) iter.Seq2[UserConsensus, error] {
	return func(yield func(UserConsensus, error) bool) {
		for remoteConsensus, err := range seq {
			consensus := parseUserConsensus(remoteConsensus)
			if !yield(consensus, err) {
				break
			}
		}
	}
}
