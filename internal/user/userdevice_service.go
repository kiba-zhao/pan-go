package user

import (
	"context"
	"encoding/binary"
	"errors"
	"iter"
	"pan/pkg/net"
	"slices"
)

var ErrUserDeviceInvalid = errors.New("user.UserDeviceService verifyUserDevice Error: UserDevice Invalid")
var ErrUserDeviceInvalidSignature = errors.New("user.UserDeviceService verifyUserDevice Error: UserDevice Invalid Signature")

type UserDeviceService struct {
	UserRepository       UserRepository
	UserDeviceRepository UserDeviceRepository

	UserDataService *UserDataService
}

func (service *UserDeviceService) ScanWithUserMetaForTopic(ctx context.Context, peerId net.PeerID, meta UserMeta) (iter.Seq2[UserDevice, error], error) {

	user, err := service.UserDataService.CheckWithUserMetaForTopic(ctx, peerId, meta)
	if err != nil {
		return nil, err
	}

	return service.UserDeviceRepository.ScanWithUserID(ctx, user.ID)
}

func (service *UserDeviceService) ScanWithUser(ctx context.Context, user User) (iter.Seq2[UserDevice, error], error) {
	return service.UserDeviceRepository.ScanWithUserID(ctx, user.ID)
}

func (service *UserDeviceService) SelectWithUser(ctx context.Context, user User, peerId net.PeerID) (UserDevice, error) {
	peerIdStr := net.EncodePeerID(peerId)
	device, err := service.UserDeviceRepository.Select(ctx, user.ID, peerIdStr)
	return device, err
}

func verifyUserPeerSignature(code string, genesisSignature string, peerId []byte, signature []byte) error {
	genesisSignatureBytes, err := DecodeSignature(genesisSignature)
	if err != nil {
		return err
	}
	signatureData := slices.Concat([]byte(code), genesisSignatureBytes, peerId)
	return net.VerifyWithPublicKeyBytes(signatureData, signature, peerId)
}

func verifyUserDevice(code string, genesisSignature string, userDevice UserDevice) ([]byte, error) {
	if len(userDevice.PeerID) == 0 {
		return nil, ErrUserDeviceInvalid
	}
	if len(userDevice.PeerSignature) == 0 {
		return nil, ErrUserDeviceInvalidSignature
	}

	peerId, err := net.DecodePeerID(userDevice.PeerID)
	if err == nil {
		err = verifyUserPeerSignature(code, genesisSignature, peerId, userDevice.PeerSignature)
	}
	if err != nil {
		return nil, err
	}

	enabled := byte(0)
	if userDevice.Enabled {
		enabled = 1
	}
	signatureData := slices.Concat(userDevice.PeerSignature, []byte(userDevice.Name), []byte(userDevice.Memo), []byte{userDevice.Level, enabled})

	return signatureData, err
}

func verifyUserDeviceHeight(code string, genesisSignature string, userDevice UserDevice) error {
	peerId, err := net.DecodePeerID(userDevice.PeerID)
	if err != nil {
		return err
	}

	genesisSignatureBytes, err := DecodeSignature(genesisSignature)
	if err != nil {
		return err
	}

	heightSignatureData := slices.Concat([]byte(code), genesisSignatureBytes)
	heightSignatureData = binary.BigEndian.AppendUint64(heightSignatureData, userDevice.Height)
	return net.VerifyWithPublicKeyBytes(heightSignatureData, userDevice.HeightSignature, peerId)
}
