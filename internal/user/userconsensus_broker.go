package user

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"iter"
	"pan/internal/feature"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

var ErrUserConsensusScanInvalidWithRemoteUser = errors.New("user.UserConsensusBroker ScanWithRemoteUserFields Error: Invalid")

type UserConsensusBroker struct {
	PeerBroker *feature.PeerBroker
}

func (broker *UserConsensusBroker) ScanWithUserMeta(peerId peer.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserConsensus, error], error) {
	ctx := context.Background()

	reader, err := broker.PeerBroker.Request(ctx, peerId, ScanUserConsensusWithUserMeta, meta)
	if err != nil {
		return nil, err
	}
	return func(yield func(*RemoteUserConsensus, error) bool) {
		defer reader.Close()

		sizeBytes := make([]byte, 2)
		var itemBytes []byte
		for {
			n, err := reader.Reader.Read(sizeBytes)
			if err == nil && n != 2 {
				err = ErrUserConsensusScanInvalidWithRemoteUser
			}
			if err == nil {
				size := binary.BigEndian.Uint16(sizeBytes)
				itemBytes = make([]byte, size)
				if size > 0 {
					n, err = reader.Reader.Read(itemBytes)
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
			err = proto.Unmarshal(itemBytes, &remoteConsensus)
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
