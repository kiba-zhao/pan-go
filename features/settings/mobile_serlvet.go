//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	libApp "pan/lib/app"
	"pan/lib/feature"
	"pan/lib/serlvet"
)

type MobileSettingsServlet struct {
	MobileSettingsService *MobileSettingsService
}

var _ = (feature.SerlvetHandler)((*MobileSettingsServlet)(nil))

func (srv *MobileSettingsServlet) SetupToSerlvet(router serlvet.SerlvetRouter) error {
	router.Handle([]byte("load"), srv.Load)
	router.Handle([]byte("save"), srv.Save)
	return nil
}

func (srv *MobileSettingsServlet) Load(ctx libApp.AppContext, next libApp.Next) error {
	settings, err := srv.MobileSettingsService.Load()
	if err != nil {
		ctx.ThrowError(serlvet.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(serlvet.CodeInternalError, err)
	}
	return nil
}

func (srv *MobileSettingsServlet) Save(ctx libApp.AppContext, next libApp.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(serlvet.CodeBadRequest, err)
		return nil
	}

	var fields MobileSettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(serlvet.CodeBadRequest, err)
		return nil
	}

	settings, err := srv.MobileSettingsService.Save(fields)
	if err != nil {
		ctx.ThrowError(serlvet.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
