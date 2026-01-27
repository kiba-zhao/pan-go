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

type UserDeviceBroker struct {
	PeerBroker *feature.PeerBroker
}

func (broker *UserDeviceBroker) ScanWithUserMeta(peerId peer.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserDevice, error], error) {
	ctx := context.Background()
	reader, err := broker.PeerBroker.Request(ctx, peerId, ScanUserDeviceWithUserMeta, meta)
	if err != nil {
		return nil, err
	}
	return func(yield func(*RemoteUserDevice, error) bool) {
		defer reader.Close()

		sizeBytes := make([]byte, 2)
		var itemBytes []byte
		for {
			n, err := reader.Read(sizeBytes)
			if err == nil && n != 2 {
				err = ErrUserConsensusScanInvalidWithRemoteUser
			}
			if err == nil {
				size := binary.BigEndian.Uint16(sizeBytes)
				itemBytes = make([]byte, size)
				if size > 0 {
					n, err = reader.Read(itemBytes)
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

			var remoteDevice RemoteUserDevice
			err = proto.Unmarshal(itemBytes, &remoteDevice)
			if !yield(&remoteDevice, err) {
				break
			}
		}
	}, nil
}

func parseUserDeviceSeq(seq iter.Seq2[*RemoteUserDevice, error]) iter.Seq2[UserDevice, error] {
	return func(yield func(UserDevice, error) bool) {
		for remoteDevice, err := range seq {
			device := parseUserDevice(remoteDevice)
			if !yield(device, err) {
				break
			}
		}
	}
}
