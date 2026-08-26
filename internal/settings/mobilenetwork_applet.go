//go:build android || ios

package settings

import (
	"bytes"
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type MobileNetworkAppletModule struct {
	MobileNetworkService *MobileNetworkService
}

var _ = (app.AppletModule)((*MobileNetworkAppletModule)(nil))

func (srv *MobileNetworkAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("mobile-network:load"), srv.Load)
	router.Handle([]byte("mobile-network:save"), srv.Save)
	return nil
}

func (srv *MobileNetworkAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	mobileNetwork, err := srv.MobileNetworkService.Load()
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}

	resp, err := json.Marshal(mobileNetwork)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	} else {
		ctx.ThrowError(app.CodeInternalError, err)
	}
	return nil
}

func (srv *MobileNetworkAppletModule) Save(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var fields MobileNetworkFields
	err = json.Unmarshal(body, &fields)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	mobileNetwork, err := srv.MobileNetworkService.Save(fields)
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}
	resp, err := json.Marshal(mobileNetwork)
	if err == nil {
		ctx.Respond(bytes.NewReader(resp))
	}
	return nil
}
