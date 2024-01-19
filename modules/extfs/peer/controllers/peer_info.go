package controllers

import (
	"bytes"
	"net/http"
	"pan/core"
	"pan/modules/extfs/services"
	"pan/peer"

	"google.golang.org/protobuf/proto"
)

// PeerInfoController
type PeerInfoController struct {
	ExtFS          *services.ExtFSService
	AuthController *AuthController
}

// Init ...
func (c *PeerInfoController) Init(app core.App[peer.Context]) {
	app.UseFn([]byte("GetPeerInfo"), c.AuthController.Auth, c.Get)
}

// Get
func (c *PeerInfoController) Get(ctx peer.Context, next core.Next) error {
	info, err := c.ExtFS.GetLatestOneToRemote()
	if err != nil {
		return ctx.ThrowError(http.StatusInternalServerError, err.Error())
	}
	bodyBytes, err := proto.Marshal(&info)
	if err != nil {
		return ctx.ThrowError(http.StatusInternalServerError, err.Error())
	}

	body := bytes.NewReader(bodyBytes)
	return ctx.Respond(body)
}
