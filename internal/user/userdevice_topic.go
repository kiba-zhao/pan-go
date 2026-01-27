package user

import (
	"encoding/binary"
	"io"
	"iter"
	"pan/internal/feature"
	"pan/internal/peer"

	"google.golang.org/protobuf/proto"
)

type UserDeviceTopic struct {
	UserDeviceService *UserDeviceService
}

func (topic *UserDeviceTopic) SetupToPeer(router peer.PeerRouter) error {
	router.Handle(ScanUserDeviceWithUserMeta, topic.ScanWithUserMeta)
	return nil
}

func (topic *UserDeviceTopic) ScanWithUserMeta(ctx peer.PeerContext, next peer.PeerNext) error {
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
	deviceSeq, err := topic.UserDeviceService.ScanWithUserMetaForTopic(peerId.(peer.PeerID), meta)
	if err == nil && deviceSeq != nil {
		deviceBytesSeq := parseUserDeviceBytesSeq(deviceSeq)
		stream := feature.NewIterStreamForSeq2(deviceBytesSeq)
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
				deviceBytes, err := proto.Marshal(remoteDevice)
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
