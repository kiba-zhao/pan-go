package user

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"iter"
	"pan/pkg/ptp"
)

var ErrUserConsensusInvalid = errors.New("user.UserConsensusService verifyUserConsensus Error: UserConsensus Invalid")
var ErrUserConsensusServiceUserNotFound = errors.New("user.UserConsensusService Error: User Not Found")
var ErrUserConsensusServiceConsensusNotFound = errors.New("user.UserConsensusService Error: Consensus Not Found")
var ErrUserConsensusServiceConsensusConflict = errors.New("user.UserConsensusService Error: Consensus Conflict")

type UserConsensusService struct {
	UserConsensusRepository UserConsensusRepository
	UserRepository          UserRepository

	UserDataService *UserDataService
}

func (service *UserConsensusService) SelectWithUser(ctx context.Context, user User) (UserConsensus, error) {
	consensus, err := service.UserConsensusRepository.Select(ctx, user.ID, user.Height)
	if err != nil {
		return consensus, err
	}

	if consensus.ID <= 0 {
		return consensus, ErrUserConsensusServiceConsensusNotFound
	}

	if consensus.Signature != user.Signature {
		return consensus, ErrUserConsensusServiceConsensusConflict
	}
	return consensus, nil
}

func (service *UserConsensusService) ScanWithUserMetaForTopic(ctx context.Context, peerId ptp.PeerID, meta UserMeta) (iter.Seq2[UserConsensus, error], error) {

	user, err := service.UserDataService.CheckWithUserMetaForTopic(ctx, peerId, meta)
	if err != nil {
		return nil, err
	}

	_, err = service.SelectWithUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return service.UserConsensusRepository.ScanWithUserIDLessThanHeight(ctx, user.ID, meta.Height)
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
		err := ptp.VerifyWithPublicKeyBytes(userConsensus.OldPassphrase, preUserConsensus.PassphraseSignature, userConsensus.UserKey)
		if err != nil {
			return err
		}

		signatureSourceHash.Write(oldPassphraseBytes)
	}
	// end verify oldPassphrase

	signatureSourceHash.Write(userConsensus.PassphraseSignature)
	signatureSourceHash.Write(userConsensus.InfoSignature)
	signatureSourceHash.Write(userConsensus.DeviceSignature)
	signatureSourceHash.Write(userConsensus.ExtraSignature)

	if len(userConsensus.PreSignature) > 0 {
		preSignatureBytes, err := DecodeSignature(userConsensus.PreSignature)
		if err != nil {
			return err
		}
		signatureSourceHash.Write(preSignatureBytes)
	}

	err = ptp.VerifyWithPublicKeyBytes(signatureSourceHash.Sum(nil), userConsensus.PeerSignature, userConsensus.PeerID)
	if err != nil {
		return err
	}

	signatureSourceHash.Write(userConsensus.PeerID)
	signatureSourceHash.Write(userConsensus.PeerSignature)
	singnatureBytes, err := DecodeSignature(userConsensus.Signature)
	if err == nil {
		err = ptp.VerifyWithPublicKeyBytes(signatureSourceHash.Sum(nil), singnatureBytes, userConsensus.UserKey)
	}
	return err
}

func generateUserInfoSignatureData(user User) []byte {
	contentHash := sha256.New()

	contentHash.Write([]byte(user.Code))
	contentHash.Write([]byte(user.Name))
	contentHash.Write([]byte(user.Memo))

	return contentHash.Sum(nil)
}

func generateUserDeviceSignatureData(code string, genesisSignature string, userDevice UserDevice) ([]byte, error) {
	return verifyUserDevice(code, genesisSignature, userDevice)
}

func generateUserExtraSignatureData(userExtra UserExtra) []byte {
	extraBytes := make([]byte, 8+len(userExtra.Signature))
	binary.BigEndian.AppendUint64(extraBytes, userExtra.SerialNumber)
	copy(extraBytes[8:], userExtra.Signature)
	return extraBytes
}
