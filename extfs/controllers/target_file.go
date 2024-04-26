package controllers

import (
	"net/http"
	"pan/app"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
)

type TargetFileController struct {
	TargetFileService *services.TargetFileService
}

func (c *TargetFileController) Init(router app.WebRouter) {
	router.GET("/target-files", c.Search)
	router.GET("/target-files/:id", c.Select)
}

func (c *TargetFileController) Search(ctx app.WebContext) {
	var conditions models.TargetFileSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := c.TargetFileService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	app.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *TargetFileController) Select(ctx app.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	targetFile, err := c.TargetFileService.Select(id)
	if err == errors.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, targetFile)
}
