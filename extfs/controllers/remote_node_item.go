package controllers

import (
	"net/http"
	"pan/app/net"
	"pan/extfs/models"
	"pan/extfs/services"
)

type RemoteNodeItemController struct {
	RemoteNodeItemService *services.RemoteNodeItemService
}

func (s *RemoteNodeItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/remote-node-items", s.Search)
	return nil
}

func (s *RemoteNodeItemController) Search(ctx net.WebContext) {
	var condition models.RemoteNodeItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := s.RemoteNodeItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
