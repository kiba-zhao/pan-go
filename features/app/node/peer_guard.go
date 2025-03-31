package node

import (
	"errors"
	appsettings "pan/features/app/settings"
	"pan/lib/peer"
)

var ErrGuardAccessRefused = errors.New("guard.PeerGuard Error: Access Refused")

type PeerGuard struct {
	AppNodeService     *AppNodeService
	AppSettingsService *appsettings.AppSettingsService
}

// Enabled checks if the peer guard is enabled
//
// It loads the current settings from the AppSettingsService and checks if the
// GuardEnabled field is true. If the field is true, the method returns true,
// otherwise it returns false.
func (g *PeerGuard) Enabled() bool {
	settings := g.AppSettingsService.Load()
	return settings.GuardEnabled
}

// Access checks if the peer is allowed to access the p2p network.
//
// It first calls the AccessWithPeerID method of the AppNodeService to check if the
// peer is allowed to access the p2p network. If the peer is not allowed to access
// the p2p network, the method returns the error.
//
// If the peer is allowed to access the p2p network, the method loads the current
// settings from the AppSettingsService and checks if the GuardAccess field is true.
// If the field is false, the method returns ErrGuardAccessRefused.
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
