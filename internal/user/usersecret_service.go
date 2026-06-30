package user

import (
	"bytes"
	"context"
	"errors"
	"pan/pkg/net"
)

var ErrUserSecretServiceUserNotFound = errors.New("user.UserSecretService Error: User Not Found")
var ErrUserSecretServiceUserConflict = errors.New("user.UserSecretService Error: User Conflict")
var ErrUserSecretServiceUserDeviceNotFound = errors.New("user.UserSecretService Error: User Device Not Found")
var ErrUserSecretServiceUserDeviceInvalid = errors.New("user.UserSecretService Error: User Device Invalid")
var ErrUserSecretServiceUserSecretConflict = errors.New("user.UserSecretService Error: User Secret Conflict")

type UserSecretService struct {
	UserRepository       UserRepository
	UserDeviceRepository UserDeviceRepository
	UserSecretRepository UserSecretRepository

	UserSecretBroker *UserSecretBroker
}

func (service *UserSecretService) Pull(ctx context.Context, peerId net.PeerID, meta UserMeta) error {

	user, err := service.UserRepository.SelectWithGenesis(ctx, meta.GenesisSignature, meta.Code)
	if err != nil {
		return err
	}
	if user.ID <= 0 {
		return ErrUserSecretServiceUserNotFound
	}
	if user.Height != meta.Height || user.Signature != meta.Signature {
		return ErrUserSecretServiceUserConflict
	}

	secret, err := service.UserSecretRepository.SelectWithUserID(ctx, user.ID)
	if err != nil {
		return err
	}
	if secret.ID > 0 && bytes.Equal(secret.UserKey, user.UserKey) {
		return nil
	}

	device, err := service.UserDeviceRepository.Select(ctx, user.ID, net.EncodePeerID(peerId))
	if err != nil {
		return err
	}
	if device.ID <= 0 {
		return ErrUserSecretServiceUserDeviceNotFound
	}

	if len(device.PeerSignature) == 0 || !device.Enabled || device.Level != DeviceLevelOwner {
		return ErrUserSecretServiceUserDeviceInvalid
	}

	remoteMeta := parseRemoteUserMeta(meta)
	remoteSecret, err := service.UserSecretBroker.Pull(ctx, peerId, remoteMeta)
	if err != nil {
		return err
	}

	secret.UserID = user.ID
	secret.UserKey = remoteSecret.UserKey
	secret.UserSecretKey = remoteSecret.UserSecretKey
	_, err = service.UserSecretRepository.SaveWithUserID(ctx, secret)
	return err
}

func (service *UserSecretService) PullForTopic(ctx context.Context, peerId net.PeerID, meta UserMeta) (UserSecret, error) {

	user, err := service.UserRepository.SelectWithGenesis(ctx, meta.GenesisSignature, meta.Code)
	if err != nil {
		return UserSecret{}, err
	}
	if user.ID <= 0 {
		return UserSecret{}, ErrUserSecretServiceUserNotFound
	}
	if user.Height != meta.Height || user.Signature != meta.Signature {
		return UserSecret{}, ErrUserSecretServiceUserConflict
	}

	device, err := service.UserDeviceRepository.Select(ctx, user.ID, net.EncodePeerID(peerId))
	if err != nil {
		return UserSecret{}, err
	}
	if device.ID <= 0 {
		return UserSecret{}, ErrUserSecretServiceUserDeviceNotFound
	}

	if len(device.PeerSignature) == 0 || !device.Enabled || device.Level != DeviceLevelOwner {
		return UserSecret{}, ErrUserSecretServiceUserDeviceInvalid
	}

	secret, err := service.UserSecretRepository.SelectWithUserID(ctx, user.ID)
	if err != nil {
		return UserSecret{}, err
	}
	if secret.ID > 0 && !bytes.Equal(secret.UserKey, user.UserKey) {
		return UserSecret{}, ErrUserSecretServiceUserSecretConflict
	}

	return secret, err
}
