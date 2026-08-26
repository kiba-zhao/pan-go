//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/pkg/web"
)

type DeviceInfoController struct {
	DeviceInfoService *DeviceInfoService
}

var _ = (web.WebController)((*DeviceInfoController)(nil))

func (ctrl *DeviceInfoController) SetupToWeb(router web.WebRouter) error {
	router.GET("/device-info", ctrl.Load)
	router.PATCH("/device-info", ctrl.Update)
	return nil
}

func (ctrl *DeviceInfoController) Load(ctx web.WebContext) {
	settings, err := ctrl.DeviceInfoService.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *DeviceInfoController) Update(ctx web.WebContext) {
	var fields DeviceInfoFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := ctrl.DeviceInfoService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
