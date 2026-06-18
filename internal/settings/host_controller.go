//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/pkg/web"
)

type HostSettingsController struct {
	HostSettingsService *HostSettingsService
}

var _ = (web.WebController)((*HostSettingsController)(nil))

func (ctrl *HostSettingsController) SetupToWeb(router web.WebRouter) error {
	router.GET("/host", ctrl.Load)
	router.PATCH("/host", ctrl.Update)
	return nil
}

func (ctrl *HostSettingsController) Load(ctx web.WebContext) {
	settings, err := ctrl.HostSettingsService.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *HostSettingsController) Update(ctx web.WebContext) {
	var fields HostSettingsFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := ctrl.HostSettingsService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
