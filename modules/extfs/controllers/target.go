package controllers

import (
	"net/http"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/web"
)

type TargetController struct {
	TargetService *services.TargetService
}

func (c *TargetController) MountWithWeb(router web.Router) {
	r := router.Group("/targets")
	r.GET("", c.Search)
}

func (c *TargetController) Search(ctx web.Context) {
	var condition models.TargetSearchCondition
	if berr := ctx.ShouldBind(&condition); berr != nil {
		ctx.AbortWithError(http.StatusBadRequest, berr)
		return
	}

	if condition.Limit <= 0 {
		condition.Limit = 100
	}

	result, err := c.TargetService.Search(&condition)
	if err == nil {
		ctx.JSON(http.StatusOK, result)
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
