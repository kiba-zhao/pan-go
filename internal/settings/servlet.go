//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	libApp "pan/internal/app"
	"pan/internal/feature"
	"pan/internal/servlet"
)

type SettingsServlet struct {
	SettingsService *SettingsService
}

var _ = (feature.ServletHandler)((*SettingsServlet)(nil))

func (srv *SettingsServlet) SetupToServlet(router servlet.ServletRouter) error {
	router.Handle([]byte("load"), srv.Load)
	router.Handle([]byte("save"), srv.Save)
	return nil
}

func (srv *SettingsServlet) Load(ctx libApp.AppContext, next libApp.Next) error {
	settings, err := srv.SettingsService.Load()
	if err != nil {
		ctx.ThrowError(servlet.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(servlet.CodeInternalError, err)
	}
	return nil
}

func (srv *SettingsServlet) Save(ctx libApp.AppContext, next libApp.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	var fields SettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	settings, err := srv.SettingsService.Save(fields)
	if err != nil {
		ctx.ThrowError(servlet.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
