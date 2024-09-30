package controllers

import (
	"net/http"
	"pan/app/net"
	"pan/extfs/models"
	"pan/extfs/services"
)

type FileItemController struct {
	FileItemService *services.FileItemService
}

func (c *FileItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/file-items", c.Search)
	return nil
}

func (c *FileItemController) Search(ctx net.WebContext) {

	var condition models.FileItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.FileItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
