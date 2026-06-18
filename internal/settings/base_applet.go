//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type BaseSettingsAppletModule struct {
	SettingsService *SettingsService
}

var _ = (app.AppletModule)((*BaseSettingsAppletModule)(nil))

func (srv *BaseSettingsAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("base:load"), srv.Load)
	router.Handle([]byte("base:save"), srv.Save)
	return nil
}

func (srv *BaseSettingsAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	settings, err := srv.SettingsService.Load()
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

func (srv *BaseSettingsAppletModule) Save(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var fields SettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	settings, err := srv.SettingsService.Save(fields)
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
