//go:build android || ios

package appinfo

import (
	"encoding/json"
	"io"
	"pan/pkg/app"
)

type QRCodeAppletModule struct {
	QRCodeService *QRCodeService
}

var _ = (app.AppletModule)((*QRCodeAppletModule)(nil))

func (ctrl *QRCodeAppletModule) SetupToApplet(router app.AppServletRouter) error {
	router.Handle([]byte("qrcode:load"), ctrl.Load)
	return nil
}

func (ctrl *QRCodeAppletModule) Load(ctx app.AppServletContext, next app.AppServletNext) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	var condition QRCodeCondition
	err = json.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(app.CodeBadRequest, err)
		return nil
	}

	reader, err := ctrl.QRCodeService.Load(condition.Size)
	if err != nil {
		ctx.ThrowError(app.CodeInternalError, err)
		return nil
	}
	ctx.Respond(reader)
	return nil
}
