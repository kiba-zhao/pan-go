package guard

import (
	"errors"
	appnode "pan/app/app_node"
	appsettings "pan/app/app_settings"
	"pan/app/peer"
)

var ErrGuardAccessRefused = errors.New("guard.PeerGuard Error: Access Refused")

type PeerGuard struct {
	AppNodeService     *appnode.AppNodeService
	AppSettingsService *appsettings.AppSettingsService
}

func (g *PeerGuard) Enabled() bool {
	settings := g.AppSettingsService.Load()
	return settings.GuardEnabled
}

func (g *PeerGuard) Access(peerId peer.PeerID) error {
	err := g.AppNodeService.AccessWithPeerID(peerId)
	if err != nil {
		return err
	}

	settings := g.AppSettingsService.Load()
	if !settings.GuardAccess {
		err = ErrGuardAccessRefused
	}
	return err
}
