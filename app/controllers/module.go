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
	router.PUT("/modules/:name/actions/set-enabled", ctrl.SetEnabled)
}

// SetEnabled, called ModuleService.SetEnabled with models.ModuleEnabled
func (ctrl *ModuleController) SetEnabled(ctx core.WebContext) {
	var enabled models.ModuleEnabled
	if err := ctx.ShouldBind(&enabled); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	name := ctx.Param("name")
	err := ctrl.ModuleService.SetEnabled(name, *enabled.Enabled)
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
	ctx.JSON(http.StatusOK, enabled)
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

	result, err := ctrl.ModuleService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
