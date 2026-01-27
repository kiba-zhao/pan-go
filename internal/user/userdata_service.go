package user

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"iter"
	"pan/internal/peer"
	"pan/internal/settings"
)

var ErrUserDataServiceConsensusNotFound = errors.New("user.UserDataService Error: Consensus Not Found")
var ErrUserDataServiceConsensusConflict = errors.New("user.UserDataService Error: Consensus Conflict")
var ErrUserDataServiceUserContentConflict = errors.New("user.UserDataService Error: User Content Conflict")
var ErrUserDataServiceDeviceContentConflict = errors.New("user.UserDataService Error: Device Content Conflict")
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

	UserDataBroker      *UserDataBroker
	UserConsensusBroker *UserConsensusBroker
	UserDeviceBroker    *UserDeviceBroker
	UserSecretBroker    *UserSecretBroker
}

func (service *UserDataService) Pull(peerId peer.PeerID, meta UserMeta) error {

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

	// Verify user content
	userContent := generateUserContent(user)
	if !bytes.Equal(userContent, prevConsensus.UserContent) {
		return ErrUserDataServiceUserContentConflict
	}
	//

	// Verify device content
	remoteDeviceSeq, err := service.UserDeviceBroker.ScanWithUserMeta(peerId, &remoteUserMeta)
	if err != nil {
		return err
	}
	deviceSeq := parseUserDeviceSeq(remoteDeviceSeq)

	var remoteDevice *UserDevice
	var hostDevice *UserDevice
	hostPeerId := service.SecurityConfig.PeerID()
	var userDevices []UserDevice
	deviceContentHash := sha256.New()
	for device, err := range deviceSeq {
		if err != nil {
			return err
		}

		// Add device content hash
		deviceContent, err := generateUserDeviceContent(user.Code, user.GenesisSignature, device)
		if err != nil {
			return err
		}
		deviceContentHash.Write(deviceContent)
		//

		devicePeerId, _ := peer.DecodePeerID(device.PeerID)
		if hostDevice == nil && bytes.Equal(devicePeerId, hostPeerId) {
			hostDevice = &device
		} else if remoteDevice == nil && bytes.Equal(devicePeerId, peerId) {
			remoteDevice = &device
		}
		userDevices = append(userDevices, device)
	}

	if !bytes.Equal(deviceContentHash.Sum(nil), prevConsensus.DeviceContent) {
		return ErrUserDataServiceDeviceContentConflict
	}

	if hostDevice == nil || remoteDevice == nil {
		return ErrUserDataServiceDeviceConflict
	}
	//

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

	return service.UserDataRepository.Save(user, secret, userConsensuses, userDevices)
}

func (service *UserDataService) PullForTopic(peerId peer.PeerID, meta UserMeta) (User, error) {
	user, err := service.UserRepository.SelectWithGenesis(meta.GenesisSignature, meta.Code)
	if err != nil {
		return User{}, err
	}
	if user.ID <= 0 {
		return User{}, ErrUserDataServiceUserNotFound
	}
	if user.Height < meta.Height {
		return User{}, ErrUserDataServiceUserConflict
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

	peerIdStr := peer.EncodePeerID(peerId)
	device, err := service.UserDeviceRepository.Select(user.ID, peerIdStr)
	if err != nil {
		return User{}, err
	}
	if device.ID <= 0 {
		return User{}, ErrUserDataServiceDeviceNotFound
	}

	return user, nil
}

func (service *UserDataService) Push(peerId peer.PeerID, meta UserMeta, userId uint) error {
	peerIdStr := peer.EncodePeerID(peerId)
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

func (service *UserDataService) PushForTopic(peerId peer.PeerID, meta UserMeta, signature []byte) error {
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
