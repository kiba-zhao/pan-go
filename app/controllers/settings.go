package controllers

import (
	"net/http"
	"pan/app/models"
	"pan/app/net"
	"pan/app/services"
)

type SettingsController struct {
	SettingsService *services.SettingsService
}

func (c *SettingsController) SetupToWeb(router net.WebRouter) error {
	router.GET("/settings", c.Load)
	router.PATCH("/settings", c.Update)
	return nil
}

func (c *SettingsController) Load(ctx net.WebContext) {
	settings := c.SettingsService.Load()
	ctx.JSON(http.StatusOK, settings)
}

func (c *SettingsController) Update(ctx net.WebContext) {
	var fields models.SettingsFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := c.SettingsService.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
