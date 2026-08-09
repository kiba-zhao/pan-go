package user

import (
	"context"
	"encoding/binary"
	"iter"
	"pan/pkg/proto"
	"pan/pkg/ptp"

	protobuf "google.golang.org/protobuf/proto"
)

type UserDeviceTopic struct {
	UserDeviceService *UserDeviceService
}

var _ = (ptp.PeerTopic)((*UserDeviceTopic)(nil))

func (topic *UserDeviceTopic) SetupToPeer(router ptp.PeerServletRouter) error {
	router.Handle(ScanUserDeviceWithUserMeta, topic.ScanWithUserMeta)
	return nil
}

func (topic *UserDeviceTopic) ScanWithUserMeta(ctx ptp.PeerServletContext, next ptp.PeerServletNext) error {
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
	deviceSeq, err := topic.UserDeviceService.ScanWithUserMetaForTopic(context.Background(), peerId.(ptp.PeerID), meta)
	if err == nil && deviceSeq != nil {
		deviceBytesSeq := parseUserDeviceBytesSeq(deviceSeq)
		stream := NewIterStreamForSeq2(deviceBytesSeq)
		ctx.Respond(stream)
	}
	return err
}

func parseUserDeviceBytesSeq(seq iter.Seq2[UserDevice, error]) iter.Seq2[[]byte, error] {
	return func(yield func([]byte, error) bool) {
		next, stop := iter.Pull2(seq)
		defer stop()
		for {
			device, err, done := next()
			var bytes []byte
			if err == nil {
				remoteDevice := parseRemoteUserDevice(device)
				deviceBytes, err := protobuf.Marshal(remoteDevice)
				if err == nil {
					size := len(deviceBytes)
					bytes = make([]byte, 2+size)
					binary.BigEndian.PutUint16(bytes, uint16(len(deviceBytes)))
					copy(bytes[2:], deviceBytes)
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
