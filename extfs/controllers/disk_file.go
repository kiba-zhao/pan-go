package controllers

import (
	"net/http"
	"pan/core"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
)

type DiskFileController struct {
	DiskFileService *services.DiskFileService
}

func (c *DiskFileController) Init(router core.WebRouter) {
	router.GET("/disk-files", c.Search)
}

func (c *DiskFileController) Search(ctx core.WebContext) {
	var conditions models.DiskFileSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := c.DiskFileService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err == errors.ErrConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	core.SetCountHeaderForWeb(ctx, total)
	if total == 0 {
		items = []models.DiskFile{}
	}
	ctx.JSON(http.StatusOK, items)
}
