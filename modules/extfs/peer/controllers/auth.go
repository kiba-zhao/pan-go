package controllers

import (
	"net/http"
	"pan/core"
	"pan/modules/extfs/services"
	"pan/peer"
)

type AuthController struct {
	PeerService *services.PeerService
}

// Auth ...
func (c *AuthController) Auth(ctx peer.Context, next core.Next) error {
	peerId := ctx.PeerId()
	hasPeer := c.PeerService.HasEnabledRemotePeer(peerId)
	if !hasPeer {
		return ctx.ThrowError(http.StatusForbidden, "Forbidden")
	}
	return next()
}
