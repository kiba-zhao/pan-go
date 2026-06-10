package user

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"iter"
	"pan/internal/net"
	"pan/internal/settings"
)

var ErrUserDataServiceConsensusNotFound = errors.New("user.UserDataService Error: Consensus Not Found")
var ErrUserDataServiceConsensusConflict = errors.New("user.UserDataService Error: Consensus Conflict")
var ErrUserDataServiceUserSignatureConflict = errors.New("user.UserDataService Error: User Signature Conflict")
var ErrUserDataServiceDeviceSignatureConflict = errors.New("user.UserDataService Error: Device Signature Conflict")
var ErrUserDataServiceExtraSignatureConflict = errors.New("user.UserDataService Error: Extra Signature Conflict")
var ErrUserDataServiceDeviceConflict = errors.New("user.UserDataService Error: Device Conflict")
var ErrUserDataServiceUserConflict = errors.New("user.UserDataService Error: User Conflict")
var ErrUserDataServiceUserNotFound = errors.New("user.UserDataService Error: User Not Found")
var ErrUserDataServiceDeviceNotFound = errors.New("user.UserDataService Error: Device Not Found")

type UserDataService struct {
	SecurityConfig          settings.SecurityConfig
	UserDataRepository      UserDataRepository
	UserRepository          UserRepository
	UserConsensusRepository UserConsensusRepository
	UserDeviceRepository    UserDeviceRepository
	UserExtraRepository     UserExtraRepository

	UserDataBroker      *UserDataBroker
	UserConsensusBroker *UserConsensusBroker
	UserDeviceBroker    *UserDeviceBroker
	UserSecretBroker    *UserSecretBroker
	UserExtraBroker     *UserExtraBroker
}

func (service *UserDataService) Pull(peerId net.PeerID, meta UserMeta) error {

	remoteMeta := parseRemoteUserMeta(meta)
	remoteUser, err := service.UserDataBroker.Pull(peerId, remoteMeta)
	if err != nil {
		return err
	}

	user := parseUser(remoteUser)
	if user.Code != meta.Code || user.GenesisSignature != meta.GenesisSignature {
		return ErrUserDataServiceUserConflict
	}

	if user.Height < meta.Height || (user.Height == meta.Height && user.Signature != meta.Signature) {
		return ErrUserDataServiceUserConflict
	}

	var remoteUserMeta RemoteUserMeta
	remoteUserMeta.Code = user.Code
	remoteUserMeta.GenesisSignature = user.GenesisSignature
	remoteUserMeta.Signature = user.Signature
	remoteUserMeta.Height = user.Height

	// Verify consensus
	remoteConsensusSeq, err := service.UserConsensusBroker.ScanWithUserMeta(peerId, &remoteUserMeta)
	if err != nil {
		return err
	}
	consensusSeq := parseUserConsensusSeq(remoteConsensusSeq)
	next, stop := iter.Pull2(consensusSeq)

	genesisConsensus, err, ok := next()
	if err != nil {
		stop()
		return err
	}
	if !ok {
		stop()
		return ErrUserDataServiceConsensusNotFound
	}

	if genesisConsensus.Height != 0 || genesisConsensus.Signature != user.GenesisSignature {
		stop()
		return ErrUserDataServiceConsensusConflict
	}

	if err = verifyUserConsensus(user.Code, user.GenesisSignature, genesisConsensus, UserConsensus{}); err != nil {
		stop()
		return err
	}
	// TODO
	userConsensuses := []UserConsensus{genesisConsensus}
	prevConsensus := genesisConsensus
	for {
		consensus, nextErr, ok := next()
		err = nextErr
		if !ok || err != nil {
			break
		}

		if err = verifyUserConsensus(user.Code, user.GenesisSignature, consensus, prevConsensus); err != nil {
			return err
		}

		userConsensuses = append(userConsensuses, consensus)
		prevConsensus = consensus
	}
	stop()
	if err != nil {
		return err
	}

	// Check latest consensus
	if prevConsensus.Height != user.Height || prevConsensus.Signature != user.Signature {
		return ErrUserDataServiceConsensusConflict
	}
	//

	// Verify user signature
	infoSignature := generateUserInfoSignatureData(user)
	if !bytes.Equal(infoSignature, prevConsensus.InfoSignature) {
		return ErrUserDataServiceUserSignatureConflict
	}
	//

	// Verify device signature
	remoteDeviceSeq, err := service.UserDeviceBroker.ScanWithUserMeta(peerId, &remoteUserMeta)
	if err != nil {
		return err
	}

	var remoteDevice *UserDevice
	var hostDevice *UserDevice
	hostPeerId := service.SecurityConfig.PeerID()
	var userDevices []UserDevice
	deviceSignatureHash := sha256.New()
	for remoteUserDevice, err := range remoteDeviceSeq {
		if err != nil {
			return err
		}
		device := parseUserDevice(remoteUserDevice)

		// Add device content hash
		deviceSignatureData, err := generateUserDeviceSignatureData(user.Code, user.GenesisSignature, device)
		if err == nil {
			err = verifyUserDeviceHeight(user.Code, user.GenesisSignature, device)
		}
		if err != nil {
			return err
		}
		deviceSignatureHash.Write(deviceSignatureData)
		//

		devicePeerId, _ := net.DecodePeerID(device.PeerID)
		if hostDevice == nil && bytes.Equal(devicePeerId, hostPeerId) {
			hostDevice = &device
		} else if remoteDevice == nil && bytes.Equal(devicePeerId, peerId) {
			remoteDevice = &device
		}
		userDevices = append(userDevices, device)
	}

	if !bytes.Equal(deviceSignatureHash.Sum(nil), prevConsensus.DeviceSignature) {
		return ErrUserDataServiceDeviceSignatureConflict
	}

	if hostDevice == nil || remoteDevice == nil {
		return ErrUserDataServiceDeviceConflict
	}
	//

	// Verify extra signature
	remoteExtraSeq, err := service.UserExtraBroker.ScanWithUserMeta(peerId, &remoteUserMeta)
	if err != nil {
		return err
	}

	var userExtras []UserExtra
	extraSignatureHash := sha256.New()
	for remoteExtra, err := range remoteExtraSeq {
		if err != nil {
			return err
		}
		extra := parseUserExtra(remoteExtra)

		// Add extra content hash
		extraSignatureData := generateUserExtraSignatureData(extra)
		extraSignatureHash.Write(extraSignatureData)
		//

		userExtras = append(userExtras, extra)
	}

	if !bytes.Equal(extraSignatureHash.Sum(nil), prevConsensus.ExtraSignature) {
		return ErrUserDataServiceExtraSignatureConflict
	}

	// End of extra signature

	var secret UserSecret
	if hostDevice.Level == DeviceLevelOwner && remoteDevice.Level == DeviceLevelOwner {
		remoteSecret, err := service.UserSecretBroker.Pull(peerId, &remoteUserMeta)
		if err != nil {
			return err
		}
		secret.UserID = user.ID
		secret.UserKey = remoteSecret.UserKey
		secret.UserSecretKey = remoteSecret.UserSecretKey
	}

	return service.UserDataRepository.Save(user, secret, userConsensuses, userDevices, userExtras)
}

func (service *UserDataService) CheckWithUserMetaForTopic(peerId net.PeerID, meta UserMeta) (User, error) {
	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return user, err
	} else if user.ID <= 0 {
		return user, ErrUserDataServiceUserNotFound
	} else if user.Height < meta.Height {
		return user, ErrUserDataServiceUserConflict
	}

	device, err := service.UserDeviceRepository.Select(user.ID, net.EncodePeerID(peerId))
	if err != nil {
		return user, err
	} else if device.ID <= 0 {
		return user, ErrUserDataServiceDeviceNotFound
	}

	return user, nil
}

func (service *UserDataService) PullForTopic(peerId net.PeerID, meta UserMeta) (User, error) {
	user, err := service.CheckWithUserMetaForTopic(peerId, meta)
	if err != nil {
		return User{}, err
	}

	consensus, err := service.UserConsensusRepository.SelectLatestWithUserID(user.ID)
	if err != nil {
		return User{}, err
	}
	if consensus.ID <= 0 || consensus.Height < meta.Height {
		return User{}, ErrUserDataServiceConsensusNotFound
	}
	if user.Height != consensus.Height || user.Signature != consensus.Signature {
		return User{}, ErrUserDataServiceConsensusConflict
	}

	if consensus.Height > meta.Height {
		consensus, err = service.UserConsensusRepository.Select(user.ID, meta.Height)
		if err != nil {
			return User{}, err
		}
		if consensus.ID <= 0 {
			return User{}, ErrUserDataServiceConsensusNotFound
		}
	}

	if meta.Height > 0 && consensus.Signature != meta.Signature {
		return User{}, ErrUserDataServiceConsensusConflict
	}

	return user, err
}

func (service *UserDataService) Push(peerId net.PeerID, meta UserMeta, userId uint) error {
	peerIdStr := net.EncodePeerID(peerId)
	device, err := service.UserDeviceRepository.Select(userId, peerIdStr)
	if err != nil {
		return err
	}
	if device.ID <= 0 {
		return ErrUserDataServiceDeviceNotFound
	}

	remoteMeta := parseRemoteUserMeta(meta)
	return service.UserDataBroker.Push(peerId, remoteMeta, device.PeerSignature)
}

func (service *UserDataService) PushForTopic(peerId net.PeerID, meta UserMeta, signature []byte) error {
	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return err
	}
	if user.ID > 0 {
		return ErrUserDataServiceUserConflict
	}

	err = verifyUserPeerSignature(meta.Code, meta.GenesisSignature, peerId, signature)
	if err != nil {
		return err
	}

	var meta_ UserMeta
	meta_.GenesisSignature = meta.GenesisSignature
	meta_.Code = meta.Code
	err = service.Pull(peerId, meta_)
	return err
}

func (service *UserDataService) Clean(meta UserMeta) error {
	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return err
	}
	if user.ID <= 0 {
		return ErrUserDataServiceUserNotFound
	}
	if user.Height != meta.Height || user.Signature != meta.Signature {
		return ErrUserDataServiceUserConflict
	}

	return service.UserDataRepository.Clean(user)
}
