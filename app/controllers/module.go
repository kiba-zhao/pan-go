package controllers

import (
	"net/http"
	"pan/app/errors"
	"pan/app/models"
	"pan/app/services"
	"pan/core"
)

type ModuleController struct {
	ModuleService *services.ModuleService
}

func (ctrl *ModuleController) Init(router core.WebRouter) {
	router.GET("/modules", ctrl.Search)
	router.GET("/modules/:name", ctrl.Get)
	router.PATCH("/modules/:name", ctrl.Update)
}

func (ctrl *ModuleController) Update(ctx core.WebContext) {
	var fields models.ModuleFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name := ctx.Param("name")
	module, err := ctrl.ModuleService.Update(name, *fields.Enabled)
	if err == errors.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err == errors.ErrForbidden {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, module)
}

func (ctrl *ModuleController) Get(ctx core.WebContext) {
	name := ctx.Param("name")
	module, err := ctrl.ModuleService.Get(name)
	if err == errors.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, module)
}

func (ctrl *ModuleController) Search(ctx core.WebContext) {
	var conditions models.ModuleSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := ctrl.ModuleService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	core.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
