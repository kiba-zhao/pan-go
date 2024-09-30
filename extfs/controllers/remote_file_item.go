package controllers

import (
	"net/http"
	"pan/app/net"
	"pan/extfs/models"
	"pan/extfs/services"
)

type RemoteFileItemController struct {
	RemoteFileItemService *services.RemoteFileItemService
}

func (c *RemoteFileItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/remote-file-items", c.Search)
	return nil
}

func (c *RemoteFileItemController) Search(ctx net.WebContext) {

	var condition models.RemoteFileItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.RemoteFileItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
