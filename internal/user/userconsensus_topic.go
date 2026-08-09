package user

import (
	"context"
	"encoding/binary"
	"iter"
	"pan/pkg/proto"
	"pan/pkg/ptp"

	protobuf "google.golang.org/protobuf/proto"
)

type UserConsensusTopic struct {
	UserConsensusService *UserConsensusService
}

var _ = (ptp.PeerTopic)((*UserConsensusTopic)(nil))

func (topic *UserConsensusTopic) SetupToPeer(router ptp.PeerServletRouter) error {
	router.Handle(ScanUserConsensusWithUserMeta, topic.ScanWithUserMeta)
	return nil
}

func (topic *UserConsensusTopic) ScanWithUserMeta(ctx ptp.PeerServletContext, next ptp.PeerServletNext) error {
	peerId, ok := ctx.Session(ptp.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(ptp.CodeBadRequest, nil)
		return nil
	}

	var userMeta RemoteUserMeta
	err := proto.UnmarshalWithReader(ctx.Request(), &userMeta)
	if err != nil {
		ctx.ThrowError(ptp.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)
	consensusSeq, err := topic.UserConsensusService.ScanWithUserMetaForTopic(context.Background(), peerId.(ptp.PeerID), meta)
	if err == nil && consensusSeq != nil {
		consensusBytesSeq := parseUserConsensusBytesSeq(consensusSeq)
		stream := NewIterStreamForSeq2(consensusBytesSeq)
		ctx.Respond(stream)
	}

	return err
}

func parseUserConsensusBytesSeq(seq iter.Seq2[UserConsensus, error]) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()
		for {
			consensus, err, done := next()
			var bytes []byte
			if err == nil {
				remoteConsensus := parseRemoteUserConsensus(consensus)
				consensusBytes, err := protobuf.Marshal(remoteConsensus)
				if err == nil {
					size := len(consensusBytes)
					bytes = make([]byte, 2+size)
					binary.BigEndian.PutUint16(bytes, uint16(len(consensusBytes)))
					copy(bytes[2:], consensusBytes)
				}
			}
			if !yield(bytes, err) {
				break
			}
			if done {
				break
			}
		}
	}
}
