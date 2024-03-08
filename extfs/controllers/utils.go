package controllers

import (
	"pan/core"
	"pan/extfs/services"
	"strconv"
)

type ListUtilsController struct {
	SettingsService *services.SettingsService
}

func (c *ListUtilsController) SetTotal(ctx core.WebContext, total int64) {
	headerName := c.SettingsService.GetTotalHeaderName()
	ctx.Header(headerName, strconv.FormatInt(total, 10))
	ctx.Header("Access-Control-Expose-Headers'", headerName)
}
