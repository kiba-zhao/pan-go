package user

import (
	"errors"
	"iter"
	"pan/internal/peer"
	"slices"
)

var ErrUserDeviceInvalid = errors.New("user.UserDeviceService verifyUserDevice Error: UserDevice Invalid")
var ErrUserDeviceInvalidSignature = errors.New("user.UserDeviceService verifyUserDevice Error: UserDevice Invalid Signature")
var ErrUserDeviceServiceUserNotFound = errors.New("user.UserDeviceService Error: User Not Found")
var ErrUserDeviceServiceConsensusNotFound = errors.New("user.UserDeviceService Error: Consensus Not Found")
var ErrUserDeviceServiceConsensusConflict = errors.New("user.UserDeviceService Error: Consensus Conflict")

type UserDeviceService struct {
	UserRepository          UserRepository
	UserConsensusRepository UserConsensusRepository
	UserDeviceRepository    UserDeviceRepository
}

func (service *UserDeviceService) ScanWithUserMetaForTopic(peerId peer.PeerID, meta UserMeta) (iter.Seq2[UserDevice, error], error) {

	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return nil, err
	}

	if user.ID <= 0 {
		return nil, ErrUserDeviceServiceUserNotFound
	}

	consensus, err := service.UserConsensusRepository.SelectLatestWithUserID(user.ID)
	if err != nil {
		return nil, err
	}

	if consensus.ID <= 0 {
		return nil, ErrUserDeviceServiceConsensusNotFound
	}

	if consensus.Signature != meta.Signature || consensus.Height != meta.Height {
		return nil, ErrUserDeviceServiceConsensusConflict
	}

	return service.UserDeviceRepository.ScanWithUserID(user.ID)
}

func verifyUserPeerSignature(code string, genesisSignature string, peerId []byte, signature []byte) error {
	genesisSignatureBytes, err := DecodeSignature(genesisSignature)
	if err != nil {
		return err
	}
	signatureData := slices.Concat([]byte(code), genesisSignatureBytes, peerId)
	return peer.VerifyWithPublicKeyBytes(signatureData, signature, peerId)
}

func verifyUserDevice(code string, genesisSignature string, userDevice UserDevice) ([]byte, error) {
	if len(userDevice.PeerID) == 0 {
		return nil, ErrUserDeviceInvalid
	}
	if len(userDevice.PeerSignature) == 0 {
		return nil, ErrUserDeviceInvalidSignature
	}

	peerId, err := peer.DecodePeerID(userDevice.PeerID)
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
	signatureData := slices.Concat(userDevice.PeerSignature, []byte{userDevice.Level, enabled})

	return signatureData, err
}
