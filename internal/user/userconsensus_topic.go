package user

import (
	"encoding/binary"
	"io"
	"iter"
	"pan/internal/feature"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

type UserConsensusTopic struct {
	UserConsensusService *UserConsensusService
}

func (topic *UserConsensusTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(ScanUserConsensusWithUserMeta, topic.ScanWithUserMeta)
	return nil
}

func (topic *UserConsensusTopic) ScanWithUserMeta(ctx peer.PeerContext, next peer.PeerNext) error {
	peerId, ok := ctx.Session(peer.ContextPeerID)
	if !ok {
		ctx.ThrowError(peer.CodeBadRequest, nil)
		return nil
	}

	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	var userMeta RemoteUserMeta
	err = proto.Unmarshal(body, &userMeta)
	if err != nil {
		ctx.ThrowError(peer.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)
	consensusSeq, err := topic.UserConsensusService.ScanWithUserMetaForTopic(peerId.(peer.PeerID), meta)
	if err == nil && consensusSeq != nil {
		consensusBytesSeq := parseUserConsensusBytesSeq(consensusSeq)
		stream := feature.NewIterStreamForSeq2(consensusBytesSeq)
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
				consensusBytes, err := proto.Marshal(remoteConsensus)
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
