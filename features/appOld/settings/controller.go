package settings

import (
	"net/http"
	"pan/lib/web"
)

type AppSettingsController struct {
	AppSettingsService *AppSettingsService
}

func (c *AppSettingsController) SetupToWeb(router web.WebRouter) error {
	router.GET("/settings", c.Load)
	router.PATCH("/settings", c.Update)
	return nil
}

func (c *AppSettingsController) Load(ctx web.WebContext) {
	settings := c.AppSettingsService.Load()
	ctx.JSON(http.StatusOK, settings)
}

func (c *AppSettingsController) Update(ctx web.WebContext) {
	var fields AppSettingsFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := c.AppSettingsService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
