package user

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"pan/pkg/log"
	"pan/pkg/net"
	"slices"
	"sync"
	"time"
)

var ErrUserDeviceSyncAgentSyncExit = errors.New("user.UserDeviceSyncAgent Error: Sync Exit")
var ErrUserDeviceSyncAgentSyncOngoing = errors.New("user.UserDeviceSyncAgent Error: Sync Ongoing")
var ErrUserDeviceSyncAgentReviseOngoing = errors.New("user.UserDeviceSyncAgent Error: Revise Ongoing")
var ErrUserDeviceSyncAgentSyncInvalidPeer = errors.New("user.UserDeviceSyncAgent Error: Sync Invalid Peer")

var UserDeviceSyncAliveTimeout = time.Second * 6

type UserDeviceSyncErr struct {
	userID      uint
	peerId      net.PeerID
	err         error
	isCompleted bool
	rw          sync.RWMutex
}

type UserDeviceSyncAgent struct {
	UserService          *UserService
	UserConsensusService *UserConsensusService
	UserDataService      *UserDataService
	UserDeviceService    *UserDeviceService

	logger log.Logger

	syncErrs []*UserDeviceSyncErr
	syncRW   sync.RWMutex

	isSyncing     bool
	syncingLocker sync.Mutex

	syncCh       chan struct{}
	syncLocker   sync.Mutex
	syncReload   bool
	syncEnabled  bool
	syncInterval time.Duration

	peerId net.PeerID
	rw     sync.RWMutex
}

func (agent *UserDeviceSyncAgent) doSync(ctx context.Context) error {
	agent.logger.Debug("user.UserDeviceSyncAgent", "sync begin")
	defer agent.logger.Debug("user.UserDeviceSyncAgent", "sync end")

	var err error
	var delayIntervalCh <-chan time.Time
	var syncEnabled bool
	var syncInterval time.Duration

	for {
		err = nil
		if delayIntervalCh == nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case <-agent.syncCh:
				agent.syncLocker.Lock()
				agent.syncReload = true
				syncEnabled = agent.syncEnabled
				syncInterval = agent.syncInterval
				agent.syncLocker.Unlock()
			}
		} else {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case <-agent.syncCh:
				agent.syncLocker.Lock()
				agent.syncReload = true
				syncEnabled = agent.syncEnabled
				syncInterval = agent.syncInterval
				agent.syncLocker.Unlock()
			case <-delayIntervalCh:
				agent.syncLocker.Lock()
				syncEnabled = agent.syncEnabled
				syncInterval = agent.syncInterval
				agent.syncLocker.Unlock()
			}
		}

		if err != nil {
			break
		}

		if !syncEnabled {
			delayIntervalCh = nil
			continue
		}

		err = agent.sync(ctx)
		if ctx.Err() != nil {
			break
		}

		delayIntervalCh = time.After(syncInterval)
	}

	return err
}

func (agent *UserDeviceSyncAgent) sync(ctx context.Context) error {

	agent.syncingLocker.Lock()
	isSyncing := agent.isSyncing
	if !isSyncing {
		agent.isSyncing = true
	}
	agent.syncingLocker.Unlock()

	if isSyncing {
		return ErrUserDeviceSyncAgentSyncOngoing
	}

	userSeq, err := agent.UserService.Scan(ctx)
	if err != nil {
		return err
	}

	for user, err := range userSeq {
		if err != nil || ctx.Err() != nil {
			break
		}

		_, err = agent.UserConsensusService.SelectWithUser(ctx, user)
		if err != nil {
			// TODO: handle error
			continue
		}

		deviceSeq, err := agent.UserDeviceService.ScanWithUser(ctx, user)
		if err != nil {
			// TODO: handle error
			continue
		}

		for device, err := range deviceSeq {
			if err != nil || ctx.Err() != nil {
				break
			}

			devicePeerId, err := net.DecodePeerID(device.PeerID)
			if err != nil {
				continue
			}

			err = agent.syncFromUserDevice(ctx, user, devicePeerId)
			// TODO: handle error
		}
	}

	agent.syncingLocker.Lock()
	agent.isSyncing = false
	agent.syncingLocker.Unlock()

	return err
}

func (agent *UserDeviceSyncAgent) syncFromUserDevice(ctx context.Context, user User, devicePeerId net.PeerID) error {

	agent.rw.RLock()
	peerId := agent.peerId
	agent.rw.RUnlock()

	if bytes.Equal(peerId, devicePeerId) {
		return ErrUserDeviceSyncAgentSyncInvalidPeer
	}

	var syncErr *UserDeviceSyncErr
	agent.syncRW.Lock()
	syncErrKey := generateUserDeviceSyncErrKey(user.ID, devicePeerId)
	syncErrIdx, syncErrOK := slices.BinarySearchFunc(agent.syncErrs, syncErrKey, compareUserDeviceSyncErr)
	if syncErrOK {
		syncErr = agent.syncErrs[syncErrIdx]
	} else {
		syncErr = &UserDeviceSyncErr{}
		syncErr.userID = user.ID
		syncErr.peerId = devicePeerId
		agent.syncErrs = slices.Insert(agent.syncErrs, syncErrIdx, syncErr)
	}
	agent.syncRW.Unlock()

	var err error
	if syncErrOK {
		syncErr.rw.RLock()
		err = syncErr.err
		isCompleted := syncErr.isCompleted
		syncErr.rw.RUnlock()

		if !isCompleted {
			return errors.New(fmt.Sprintf("user.UserDeviceSyncAgent Error: Sync Ongoing %s", net.EncodePeerID(devicePeerId)))
		}
	}

	meta := parseUserMetaWithUser(user)
	err = agent.UserDataService.Pull(ctx, devicePeerId, meta)
	if err != nil {
		errStr := err.Error()
		if errStr == ErrUserDataServiceUserNotFound.Error() || errStr == ErrUserDataServiceUserConflict.Error() {
			err = agent.UserDataService.Push(ctx, devicePeerId, meta, user.ID)
		}
	}

	syncErr.rw.Lock()
	syncErr.err = err
	syncErr.isCompleted = true
	syncErr.rw.Unlock()

	return err
}

func (agent *UserDeviceSyncAgent) purge(user User, peerId net.PeerID) {
	agent.syncRW.Lock()
	defer agent.syncRW.Unlock()
	syncErrKey := generateUserDeviceSyncErrKey(user.ID, peerId)
	syncErrIdx, syncErrOK := slices.BinarySearchFunc(agent.syncErrs, syncErrKey, compareUserDeviceSyncErr)
	if syncErrOK {
		agent.syncErrs = slices.Delete(agent.syncErrs, syncErrIdx, syncErrIdx+1)
	}
}

func (agent *UserDeviceSyncAgent) update(ctx context.Context, peerId net.PeerID) error {
	userSeq, err := agent.UserService.Scan(ctx)
	if err != nil {
		return err
	}

	for user, err := range userSeq {
		if err != nil || ctx.Err() != nil {
			break
		}

		_, err = agent.UserConsensusService.SelectWithUser(ctx, user)
		if err != nil {
			// TODO: handle error
			continue
		}

		_, err = agent.UserDeviceService.SelectWithUser(ctx, user, peerId)
		if err != nil {
			// TODO: handle error
			continue
		}

		err = agent.syncFromUserDevice(ctx, user, peerId)
	}
	return nil
}

func generateUserDeviceSyncErrKey(userID uint, peerId net.PeerID) []byte {
	key := make([]byte, 4+len(peerId))
	binary.BigEndian.PutUint32(key, uint32(userID))
	copy(key[4:], peerId)
	return key
}

func compareUserDeviceSyncErr(err *UserDeviceSyncErr, key []byte) int {
	errKey := generateUserDeviceSyncErrKey(err.userID, err.peerId)
	return bytes.Compare(errKey, key)
}
