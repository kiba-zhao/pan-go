//go:build !(android || ios)

package settings

import (
	"net/http"
	"pan/pkg/web"
)

type WebHostController struct {
	Service *WebHostService
}

var _ = (web.WebController)((*WebHostController)(nil))

func (ctrl *WebHostController) SetupToWeb(router web.WebRouter) error {
	router.GET("/web-host", ctrl.Load)
	router.PATCH("/web-host", ctrl.Update)
	return nil
}

func (ctrl *WebHostController) Load(ctx web.WebContext) {
	settings, err := ctrl.Service.Load()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}

func (ctrl *WebHostController) Update(ctx web.WebContext) {
	var fields WebHostFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	settings, err := ctrl.Service.Save(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, settings)
}
