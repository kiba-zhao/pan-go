//go:build android || ios

package appinfo

import (
	"encoding/json"
	"io"
	libApp "pan/internal/app"
	"pan/internal/feature"
	"pan/internal/servlet"
)

type QRCodeServlet struct {
	QRCodeService *QRCodeService
}

var _ = (feature.ServletHandler)((*QRCodeServlet)(nil))

func (srv *QRCodeServlet) SetupToServlet(router servlet.ServletRouter) error {
	router.Handle([]byte("qrcode:load"), srv.Load)
	return nil
}

func (srv *QRCodeServlet) Load(ctx libApp.AppContext, next libApp.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	var condition QRCodeCondition
	err = json.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(servlet.CodeBadRequest, err)
		return nil
	}

	reader, err := srv.QRCodeService.Load(condition.Size)
	if err != nil {
		ctx.ThrowError(servlet.CodeInternalError, err)
		return nil
	}
	ctx.Respond(reader)
	return nil
}
