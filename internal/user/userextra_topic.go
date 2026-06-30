package user

import (
	"context"
	"encoding/binary"
	"iter"
	"pan/pkg/net"
	"pan/pkg/proto"

	protobuf "google.golang.org/protobuf/proto"
)

type UserExtraTopic struct {
	UserExtraService *UserExtraService
}

var _ = (net.PeerTopic)((*UserExtraTopic)(nil))

func (topic *UserExtraTopic) SetupToPeer(router net.PeerServletRouter) error {
	router.Handle(ScanUserExtraWithUserMeta, topic.ScanWithUserMeta)
	return nil
}

func (topic *UserExtraTopic) ScanWithUserMeta(ctx net.PeerServletContext, next net.PeerServletNext) error {
	peerId, ok := ctx.Session(net.PeerIDSessionKey)
	if !ok {
		ctx.ThrowError(net.CodeBadRequest, nil)
		return nil
	}

	var userMeta RemoteUserMeta
	err := proto.UnmarshalWithReader(ctx.Request(), &userMeta)
	if err != nil {
		ctx.ThrowError(net.CodeBadRequest, err)
		return nil
	}

	meta := parseUserMeta(&userMeta)

	extraSeq, err := topic.UserExtraService.ScanWithUserMeta(context.Background(), peerId.(net.PeerID), meta)
	if err == nil && extraSeq != nil {
		extraBytesSeq := parseUserExtraBytesSeq(extraSeq)
		stream := NewIterStreamForSeq2(extraBytesSeq)
		ctx.Respond(stream)
	}

	return err
}

func parseUserExtraBytesSeq(extraSeq iter.Seq2[UserExtra, error]) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		next, stop := iter.Pull2(extraSeq)
		defer stop()
		for {
			extra, err, done := next()
			var bytes []byte
			if err == nil {
				remoteExtra := parseRemoteUserExtra(extra)
				extraBytes, err := protobuf.Marshal(remoteExtra)
				if err == nil {
					size := len(extraBytes)
					bytes = make([]byte, 2+size)
					binary.BigEndian.PutUint16(bytes, uint16(len(extraBytes)))
					copy(bytes[2:], extraBytes)
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
