//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type DeviceInfoAppletModule struct {
	DeviceInfoService *DeviceInfoService
}

var _ = (app.AppletModule)((*DeviceInfoAppletModule)(nil))

func (srv *DeviceInfoAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("device-info:load"), srv.Load)
	router.Handle([]byte("device-info:save"), srv.Save)
	return nil
}

func (srv *DeviceInfoAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	deviceInfo, err := srv.DeviceInfoService.Load()
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(deviceInfo)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(app.CodeInternalError, err)
	}
	return nil
}

func (srv *DeviceInfoAppletModule) Save(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var fields DeviceInfoFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	deviceInfo, err := srv.DeviceInfoService.Save(fields)
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(deviceInfo)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
