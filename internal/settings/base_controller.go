//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/pkg/web"
)

type BaseSettingsController struct {
	SettingsService *SettingsService
}

var _ = (web.WebController)((*BaseSettingsController)(nil))

func (ctrl *BaseSettingsController) SetupToWeb(router web.WebRouter) error {
	router.GET("/base", ctrl.Load)
	router.PATCH("/base", ctrl.Update)
	return nil
}

func (ctrl *BaseSettingsController) Load(ctx web.WebContext) {
	settings, err := ctrl.SettingsService.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *BaseSettingsController) Update(ctx web.WebContext) {
	var fields SettingsFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := ctrl.SettingsService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
