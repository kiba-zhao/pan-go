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

type MobileSettingsServlet struct {
	MobileSettingsService *MobileSettingsService
}

var _ = (feature.ServletHandler)((*MobileSettingsServlet)(nil))

func (srv *MobileSettingsServlet) SetupToServlet(router servlet.ServletRouter) error {
	router.Handle([]byte("load"), srv.Load)
	router.Handle([]byte("save"), srv.Save)
	return nil
}

func (srv *MobileSettingsServlet) Load(ctx libApp.AppContext, next libApp.Next) error {
	settings, err := srv.MobileSettingsService.Load()
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

func (srv *MobileSettingsServlet) Save(ctx libApp.AppContext, next libApp.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	var fields MobileSettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	settings, err := srv.MobileSettingsService.Save(fields)
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
