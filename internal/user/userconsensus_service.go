package user

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"iter"
	"pan/internal/peer"
)

var ErrUserConsensusInvalid = errors.New("user.UserConsensusService verifyUserConsensus Error: UserConsensus Invalid")
var ErrUserConsensusServiceUserNotFound = errors.New("user.UserConsensusService Error: User Not Found")
var ErrUserConsensusServiceConsensusNotFound = errors.New("user.UserConsensusService Error: Consensus Not Found")
var ErrUserConsensusServiceConsensusConflict = errors.New("user.UserConsensusService Error: Consensus Conflict")

type UserConsensusService struct {
	UserConsensusRepository UserConsensusRepository
	UserRepository          UserRepository
}

func (service *UserConsensusService) ScanWithUserMetaForTopic(peerId peer.PeerID, meta UserMeta) (iter.Seq2[UserConsensus, error], error) {

	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return nil, err
	}

	if user.ID <= 0 {
		return nil, ErrUserConsensusServiceUserNotFound
	}

	consensus, err := service.UserConsensusRepository.Select(user.ID, meta.Height)
	if err != nil {
		return nil, err
	}

	if consensus.ID <= 0 {
		return nil, ErrUserConsensusServiceConsensusNotFound
	}

	if consensus.Signature != meta.Signature {
		return nil, ErrUserConsensusServiceConsensusConflict
	}

	return service.UserConsensusRepository.ScanWithUserIDLessThanHeight(user.ID, meta.Height)
}

func verifyUserConsensus(code string, genesisSignature string, userConsensus UserConsensus, preUserConsensus UserConsensus) error {

	if len(userConsensus.Signature) <= 0 {
		return ErrUserConsensusInvalid
	}

	if len(userConsensus.PreSignature) > 0 {
		if preUserConsensus.Signature != userConsensus.PreSignature {
			return ErrUserConsensusInvalid
		}
		if len(userConsensus.OldPassphrase) <= 0 && !bytes.Equal(userConsensus.PassphraseSignature, preUserConsensus.PassphraseSignature) {
			return ErrUserConsensusInvalid
		}
		if userConsensus.Height <= preUserConsensus.Height {
			return ErrUserConsensusInvalid
		}
	}

	if len(userConsensus.OldPassphrase) > 0 && len(preUserConsensus.PassphraseSignature) <= 0 {
		return ErrUserConsensusInvalid
	}

	genesisSignatureBytes, err := DecodeSignature(genesisSignature)
	if err != nil {
		return err
	}

	signatureSourceHash := sha512.New()
	signatureSourceHash.Write([]byte(code))
	signatureSourceHash.Write(genesisSignatureBytes)
	signatureSourceHash.Write(binary.BigEndian.AppendUint64(nil, userConsensus.Height))
	signatureSourceHash.Write(userConsensus.UserKey)

	signatureSourceHash.Write(binary.BigEndian.AppendUint64(nil, userConsensus.Timestamp))
	// verify oldPassphrase
	var oldPassphraseBytes []byte
	if len(userConsensus.OldPassphrase) > 0 {
		err := peer.VerifyWithPublicKeyBytes(userConsensus.OldPassphrase, preUserConsensus.PassphraseSignature, userConsensus.UserKey)
		if err != nil {
			return err
		}

		signatureSourceHash.Write(oldPassphraseBytes)
	}
	// end verify oldPassphrase

	signatureSourceHash.Write(userConsensus.PassphraseSignature)
	signatureSourceHash.Write(userConsensus.UserContent)
	signatureSourceHash.Write(userConsensus.DeviceContent)

	if len(userConsensus.PreSignature) > 0 {
		preSignatureBytes, err := DecodeSignature(userConsensus.PreSignature)
		if err != nil {
			return err
		}
		signatureSourceHash.Write(preSignatureBytes)
	}

	err = peer.VerifyWithPublicKeyBytes(signatureSourceHash.Sum(nil), userConsensus.PeerSignature, userConsensus.PeerID)
	if err != nil {
		return err
	}

	signatureSourceHash.Write(userConsensus.PeerID)
	signatureSourceHash.Write(userConsensus.PeerSignature)
	singnatureBytes, err := DecodeSignature(userConsensus.Signature)
	if err == nil {
		err = peer.VerifyWithPublicKeyBytes(signatureSourceHash.Sum(nil), singnatureBytes, userConsensus.UserKey)
	}
	return err
}

func generateUserContent(user User) []byte {
	contentHash := sha256.New()

	contentHash.Write([]byte(user.Code))
	contentHash.Write([]byte(user.Name))
	contentHash.Write([]byte(user.Memo))

	return contentHash.Sum(nil)
}

func generateUserDeviceContent(code string, genesisSignature string, userDevice UserDevice) ([]byte, error) {
	return verifyUserDevice(code, genesisSignature, userDevice)
}
