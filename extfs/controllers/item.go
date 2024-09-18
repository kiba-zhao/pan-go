package controllers

import (
	"net/http"
	"pan/app/net"
	"pan/extfs/models"
	"pan/extfs/services"
)

type ItemController struct {
	ItemService *services.ItemService
}

func (c *ItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/items", c.Search)
	return nil
}

func (c *ItemController) Search(ctx net.WebContext) {
	var conditions models.ItemSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.ItemService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
