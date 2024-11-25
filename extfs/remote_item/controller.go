package remoteitem

import (
	"net/http"
	"pan/app/web"
)

type RemoteItemController struct {
	RemoteItemService *RemoteItemService
}

func (ctrl *RemoteItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote-items", ctrl.Search)
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
