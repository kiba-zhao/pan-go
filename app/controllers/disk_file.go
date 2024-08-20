package controllers

import (
	"net/http"
	"pan/app/constant"
	"pan/app/models"
	"pan/app/net"
	"pan/app/services"
)

type DiskFileController struct {
	DiskFileService *services.DiskFileService
}

func (c *DiskFileController) SetupToWeb(router net.WebRouter) error {
	router.GET("/disk-files", c.Search)
	return nil
}

func (c *DiskFileController) Search(ctx net.WebContext) {
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
	if err == constant.ErrConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	if total == 0 {
		items = []models.DiskFile{}
	}
	ctx.JSON(http.StatusOK, items)
}
