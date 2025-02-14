package remoteitem

import (
	"errors"
	"net/http"
	"pan/app/web"
	nodeitem "pan/extfs/node_item"
)

type RemoteItemController struct {
	RemoteItemService *RemoteItemService
}

func (ctrl *RemoteItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remotes/:peerId/remote-items", ctrl.Search)
	router.GET("/remotes/:peerId/remote-items/:id", ctrl.Select)
	return nil
}

func (ctrl *RemoteItemController) Search(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := ctrl.RemoteItemService.Search(peerId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *RemoteItemController) Select(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := nodeitem.ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	remoteItem, err := c.RemoteItemService.Select(peerId, id)
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

func ExtractPeerIdWithParam(name string, ctx web.WebContext) (string, error) {
	peerId := ctx.Param("peerId")
	if len(peerId) <= 0 {
		return "", errors.New("invalid peer id")
	}
	return peerId, nil
}
