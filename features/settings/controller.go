//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/lib/feature"
	"pan/lib/web"
)

type SettingsController struct {
	SettingsService *SettingsService
}

var _ = (feature.WebController)((*SettingsController)(nil))

func (ctrl *SettingsController) SetupToWeb(router web.WebRouter) error {
	router.GET("/", ctrl.Load)
	router.PATCH("/", ctrl.Update)
	return nil
}

func (ctrl *SettingsController) Load(ctx web.WebContext) {
	settings, err := ctrl.SettingsService.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *SettingsController) Update(ctx web.WebContext) {
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
