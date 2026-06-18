//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type MobileSettingsAppletModule struct {
	MobileSettingsService *MobileSettingsService
}

var _ = (app.AppletModule)((*MobileSettingsAppletModule)(nil))

func (srv *MobileSettingsAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("mobile:load"), srv.Load)
	router.Handle([]byte("mobile:save"), srv.Save)
	return nil
}

func (srv *MobileSettingsAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	settings, err := srv.MobileSettingsService.Load()
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(app.CodeInternalError, err)
	}
	return nil
}

func (srv *MobileSettingsAppletModule) Save(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var fields MobileSettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	settings, err := srv.MobileSettingsService.Save(fields)
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
