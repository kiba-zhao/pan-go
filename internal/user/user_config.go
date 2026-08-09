package user

import (
	"pan/internal/settings"
	"pan/pkg/ptp"
	"time"
)

var UserConfigSyncInterval = time.Second * 6

type UserConfig struct {
	peerId       ptp.PeerID
	syncInterval time.Duration
}

func NewUserConfig(security settings.SecurityConfig) *UserConfig {
	cfg := &UserConfig{}
	cfg.peerId = security.PeerID()
	cfg.syncInterval = UserConfigSyncInterval
	return cfg
}

func (cfg *UserConfig) PeerID() ptp.PeerID {
	return cfg.peerId
}

func (cfg *UserConfig) SyncInterval() time.Duration {
	return cfg.syncInterval
}
