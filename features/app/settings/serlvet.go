package settings

import (
	"bytes"
	"encoding/json"
	"io"
	libApp "pan/lib/app"
	"pan/lib/feature"
	"pan/lib/serlvet"
)

type AppSettingsSerlvet struct {
	AppSettingsService *AppSettingsService
}

var _ = (feature.SerlvetHandler)((*AppSettingsSerlvet)(nil))

func (s *AppSettingsSerlvet) SetupToSerlvet(router serlvet.SerlvetRouter) error {
	router.Handle([]byte("settings/load"), s.Load)
	router.Handle([]byte("settings/save"), s.Save)
	return nil
}

func (s *AppSettingsSerlvet) Load(ctx libApp.AppContext, next libApp.Next) error {
	settings := s.AppSettingsService.Load()
	resp, err := json.Marshal(settings)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(serlvet.CodeInternalError, err)
	}
	return nil
}

func (s *AppSettingsSerlvet) Save(ctx libApp.AppContext, next libApp.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(serlvet.CodeBadRequest, err)
		return nil
	}

	var fields AppSettingsFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(serlvet.CodeBadRequest, err)
		return nil
	}

	settings, err := s.AppSettingsService.Save(fields)
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
