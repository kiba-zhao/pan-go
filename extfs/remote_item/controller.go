package remoteitem

import (
	"errors"
	"net/http"
	"pan/app/web"
)

type RemoteItemController struct {
	RemoteItemService *RemoteItemService
}

func (ctrl *RemoteItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote-items", ctrl.Search)
	router.GET("/remote-items/:id", ctrl.Select)
	return nil
}

func (ctrl *RemoteItemController) Search(ctx web.WebContext) {
	var condition RemoteItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := ctrl.RemoteItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *RemoteItemController) Select(ctx web.WebContext) {
	paramId := ctx.Param("id")
	remoteItem, err := c.RemoteItemService.Select(paramId)
	if errors.Is(err, ErrRemoteItemInvalidID) {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.RemoteItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, remoteItem)
}
