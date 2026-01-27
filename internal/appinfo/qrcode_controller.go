//go:build !(android || ios)

package appinfo

import (
	"net/http"
	"pan/internal/feature"
	"pan/internal/web"
)

type QRCodeController struct {
	QRCodeService *QRCodeService
}

var _ = (feature.WebController)((*QRCodeController)(nil))

func (ctrl *QRCodeController) SetupToWeb(router web.WebRouter) error {
	router.GET("/qrcode", ctrl.Load)
	return nil
}

func (ctrl *QRCodeController) Load(ctx web.WebContext) {
	var condition QRCodeCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reader, err := ctrl.QRCodeService.Load(condition.Size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	http.ServeContent(ctx.Writer, ctx.Request, reader.Name(), reader.ModTime(), reader)
}
