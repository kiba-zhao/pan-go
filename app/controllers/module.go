package controllers

import (
	"net/http"
	"pan/app/models"
	"pan/app/services"
	"pan/core"
)

type ModuleController struct {
	ModuleService *services.ModuleService
}

func (ctrl *ModuleController) Init(router core.WebRouter) {
	router.GET("/modules", ctrl.Search)
}

func (ctrl *ModuleController) Search(ctx core.WebContext) {
	var conditions models.ModuleSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := ctrl.ModuleService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
