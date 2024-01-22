package controllers

import (
	"bytes"
	"net/http"
	"pan/core"
	"pan/modules/extfs/services"
	"pan/peer"

	"google.golang.org/protobuf/proto"
)

// RemoteFilesStateController
type RemoteFilesStateController struct {
	FilesStateService *services.FilesStateService
	GuardController   *GuardController
}

// Init ...
func (c *RemoteFilesStateController) Init(app core.App[peer.Context]) {
	app.UseFn([]byte("GetRemoteFilesState"), c.GuardController.Auth, c.Get)
}

// Get
func (c *RemoteFilesStateController) Get(ctx peer.Context, next core.Next) error {
	info, err := c.FilesStateService.GetLastOneToRemote()
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
