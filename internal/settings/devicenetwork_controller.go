//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/pkg/web"
)

type DeviceNetworkController struct {
	DeviceNetworkService *DeviceNetworkService
}

var _ = (web.WebController)((*DeviceNetworkController)(nil))

func (ctrl *DeviceNetworkController) SetupToWeb(router web.WebRouter) error {
	router.GET("/device-network", ctrl.Load)
	router.PATCH("/device-network", ctrl.Update)
	return nil
}

func (ctrl *DeviceNetworkController) Load(ctx web.WebContext) {
	settings, err := ctrl.DeviceNetworkService.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *DeviceNetworkController) Update(ctx web.WebContext) {
	var fields DeviceNetworkFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := ctrl.DeviceNetworkService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
