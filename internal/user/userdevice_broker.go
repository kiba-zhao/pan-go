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

var ErrUserDeviceScanInvalidWithRemoteUser = errors.New("user.UserDeviceBroker ScanWithUserMeta Error: Invalid")

type UserDeviceBroker struct {
	PeerBroker *net.PeerBroker
}

func (broker *UserDeviceBroker) ScanWithUserMeta(ctx context.Context, peerId net.PeerID, meta *RemoteUserMeta) (iter.Seq2[*RemoteUserDevice, error], error) {
	reader, err := proto.MarshalWithReader(meta)
	if err != nil {
		return nil, err
	}

	req := broker.PeerBroker.NewRequest(ScanUserDeviceWithUserMeta, reader)
	_, res, err := broker.PeerBroker.DoAction(ctx, peerId, req)
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
				err = ErrUserDeviceScanInvalidWithRemoteUser
			}
			if err == nil {
				size := binary.BigEndian.Uint16(sizeBytes)
				itemBytes = make([]byte, size)
				if size > 0 {
					n, err = res.Read(itemBytes)
					if err == nil && n != int(size) {
						err = ErrUserDeviceScanInvalidWithRemoteUser
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
