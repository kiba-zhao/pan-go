package controllers

import (
	"net/http"
	"pan/core"
	"pan/modules/extfs/services"
	"pan/peer"
)

type GuardController struct {
	RemotePeerService *services.RemotePeerService
}

// Auth ...
func (c *GuardController) Auth(ctx peer.Context, next core.Next) error {
	peerId := ctx.PeerId()
	hasEnabled := c.RemotePeerService.HasEnabled(peerId)
	if !hasEnabled {
		return ctx.ThrowError(http.StatusForbidden, "Forbidden")
	}
	return next()
}
