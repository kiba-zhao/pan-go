//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type DeviceNetworkAppletModule struct {
	DeviceNetworkService *DeviceNetworkService
}

var _ = (app.AppletModule)((*DeviceNetworkAppletModule)(nil))

func (srv *DeviceNetworkAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("device-network:load"), srv.Load)
	router.Handle([]byte("device-network:save"), srv.Save)
	return nil
}

func (srv *DeviceNetworkAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	deviceNetwork, err := srv.DeviceNetworkService.Load()
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(deviceNetwork)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(app.CodeInternalError, err)
	}
	return nil
}

func (srv *DeviceNetworkAppletModule) Save(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var fields DeviceNetworkFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	deviceNetwork, err := srv.DeviceNetworkService.Save(fields)
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(deviceNetwork)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
