package controllers

import (
	"net/http"
	"pan/core"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
)

type TargetController struct {
	TargetService *services.TargetService
}

func (c *TargetController) Init(router core.WebRouter) {
	router.GET("/targets", c.Search)
}

func (c *TargetController) Search(ctx core.WebContext) {
	var conditions models.TargetSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := c.TargetService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Header("Content-Range", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, items)
}
