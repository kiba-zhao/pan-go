package searchitem

import (
	"net/http"
	"pan/app/web"
	"strconv"
)

type SearchItemController struct {
	SearchItemService *SearchItemService
}

func (ctrl *SearchItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/search-items", ctrl.Search)
	router.DELETE("/search-items/:id", ctrl.Delete)
	return nil
}

func (ctrl *SearchItemController) Search(ctx web.WebContext) {
	var condition SearchItemCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := ctrl.SearchItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (ctrl *SearchItemController) Delete(ctx web.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = ctrl.SearchItemService.Delete(uint(id))
	if ctrl.SearchItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
