package user

var PullUserData = []byte("pull_user_data")
var PushUserData = []byte("push_user_data")
var ScanUserConsensusWithUserMeta = []byte("scan_user_consensus_with_user_meta")
var ScanUserDeviceWithUserMeta = []byte("scan_user_device_with_user_meta")
var ScanUserExtraWithUserMeta = []byte("scan_user_extra_with_user_meta")
var PullUserSecret = []byte("pull_user_secret")

var PeerSignatureHeaderName = []byte("peer_signature")

func parseRemoteUserMeta(meta UserMeta) *RemoteUserMeta {
	var remoteMeta RemoteUserMeta
	remoteMeta.Code = meta.Code
	remoteMeta.GenesisSignature = meta.GenesisSignature
	remoteMeta.Signature = meta.Signature
	remoteMeta.Height = meta.Height
	return &remoteMeta
}

func parseUserMeta(meta *RemoteUserMeta) UserMeta {
	var userMeta UserMeta
	userMeta.Code = meta.Code
	userMeta.GenesisSignature = meta.GenesisSignature
	userMeta.Signature = meta.Signature
	userMeta.Height = meta.Height
	return userMeta
}

func parseUser(remoteUser *RemoteUser) User {
	var user User
	if remoteUser == nil {
		return user
	}
	user.Code = remoteUser.Code
	user.GenesisSignature = remoteUser.GenesisSignature
	user.Signature = remoteUser.Signature
	user.Height = remoteUser.Height
	user.UserKey = remoteUser.UserKey
	user.Name = remoteUser.Name
	user.Memo = remoteUser.Memo
	return user
}

func parseRemoteUser(remoteUser User) *RemoteUser {
	var user RemoteUser
	user.Code = remoteUser.Code
	user.GenesisSignature = remoteUser.GenesisSignature
	user.Signature = remoteUser.Signature
	user.Height = remoteUser.Height
	user.UserKey = remoteUser.UserKey
	user.Name = remoteUser.Name
	user.Memo = remoteUser.Memo
	return &user
}

func parseUserConsensus(consensus *RemoteUserConsensus) UserConsensus {
	var userConsensus UserConsensus
	if consensus == nil {
		return userConsensus
	}
	userConsensus.Signature = consensus.Signature
	userConsensus.PreSignature = consensus.PreSignature

	userConsensus.PeerID = consensus.PeerID
	userConsensus.PeerSignature = consensus.PeerSignature

	userConsensus.Timestamp = consensus.Timestamp

	userConsensus.UserKey = consensus.UserKey
	userConsensus.Height = consensus.Height

	userConsensus.PassphraseSignature = consensus.PassphraseSignature
	userConsensus.OldPassphrase = consensus.OldPassphrase
	userConsensus.InfoSignature = consensus.InfoSignature
	userConsensus.DeviceSignature = consensus.DeviceSignature
	userConsensus.ExtraSignature = consensus.ExtraSignature

	return userConsensus
}

func parseRemoteUserConsensus(userConsensus UserConsensus) *RemoteUserConsensus {
	var consensus RemoteUserConsensus
	consensus.Signature = userConsensus.Signature
	consensus.PreSignature = userConsensus.PreSignature

	consensus.PeerID = userConsensus.PeerID
	consensus.PeerSignature = userConsensus.PeerSignature

	consensus.Timestamp = userConsensus.Timestamp

	consensus.UserKey = userConsensus.UserKey
	consensus.Height = userConsensus.Height

	consensus.PassphraseSignature = userConsensus.PassphraseSignature
	consensus.OldPassphrase = userConsensus.OldPassphrase
	consensus.InfoSignature = userConsensus.InfoSignature
	consensus.DeviceSignature = userConsensus.DeviceSignature
	consensus.ExtraSignature = userConsensus.ExtraSignature

	return &consensus
}

func parseUserDevice(device *RemoteUserDevice) UserDevice {
	var userDevice UserDevice
	if device == nil {
		return userDevice
	}
	userDevice.PeerID = device.PeerID
	userDevice.Level = device.Level[0]
	userDevice.PeerSignature = device.PeerSignature
	userDevice.Enabled = device.Enabled

	userDevice.Height = device.Height
	userDevice.HeightSignature = device.HeightSignature
	return userDevice
}

func parseRemoteUserDevice(userDevice UserDevice) *RemoteUserDevice {
	var device RemoteUserDevice
	device.PeerID = userDevice.PeerID
	device.Level = []byte{userDevice.Level}
	device.PeerSignature = userDevice.PeerSignature
	device.Enabled = userDevice.Enabled

	device.Height = userDevice.Height
	device.HeightSignature = userDevice.HeightSignature
	return &device
}

func parseUserExtra(extra *RemoteUserExtra) UserExtra {
	var userExtra UserExtra
	if extra == nil {
		return userExtra
	}
	userExtra.SerialNumber = extra.SerialNumber
	userExtra.Signature = extra.Signature
	return userExtra
}

func parseRemoteUserExtra(extra UserExtra) *RemoteUserExtra {
	var remoteExtra RemoteUserExtra
	remoteExtra.SerialNumber = extra.SerialNumber
	remoteExtra.Signature = extra.Signature
	return &remoteExtra
}

func parseUserSecret(remoteSecret *RemoteUserSecret) UserSecret {
	var secret UserSecret
	if remoteSecret == nil {
		return secret
	}
	secret.UserKey = remoteSecret.UserKey
	secret.UserSecretKey = remoteSecret.UserSecretKey
	return secret
}

func parseRemoteUserSecret(secret UserSecret) *RemoteUserSecret {
	var remoteSecret RemoteUserSecret
	remoteSecret.UserKey = secret.UserKey
	remoteSecret.UserSecretKey = secret.UserSecretKey
	return &remoteSecret
}
