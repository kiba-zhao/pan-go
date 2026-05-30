package user

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"iter"
	"pan/internal/net"
	"pan/internal/proto"

	protobuf "google.golang.org/protobuf/proto"
)

type UserDeviceBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserDeviceBroker) ScanWithUserMeta(peerId net.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserDevice, error], error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(ScanUserDeviceWithUserMeta, reader)
	_, res, err := broker.PeerBroker.DoAction(context.Background(), peerId, req)
	if err != nil {
		if res != nil {
			defer res.Close()
		}
		return nil, err
	}

	return func(yield func(*RemoteUserDevice, error) bool) {
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

			var remoteDevice RemoteUserDevice
			err = protobuf.Unmarshal(itemBytes, &remoteDevice)
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
