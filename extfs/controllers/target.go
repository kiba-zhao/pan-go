package controllers

import (
	"net/http"
	"pan/app"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
)

type TargetController struct {
	TargetService *services.TargetService
}

func (c *TargetController) Init(router app.WebRouter) {
	router.GET("/targets", c.Search)
	router.POST("/targets", c.Create)
	router.PATCH("/targets/:id", c.Update)
	router.GET("/targets/:id", c.Select)
	router.DELETE("/targets/:id", c.Delete)
}

func (c *TargetController) Search(ctx app.WebContext) {
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
	app.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *TargetController) Create(ctx app.WebContext) {
	var fields models.TargetFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	target, err := c.TargetService.Create(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, target)
}

func (c *TargetController) Update(ctx app.WebContext) {
	var fields models.TargetFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var opts models.TargetQueryOptions
	if err := ctx.ShouldBindQuery(&opts); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	target, err := c.TargetService.Update(fields, uint(id), opts)
	if err == errors.ErrConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, target)
}

func (c *TargetController) Select(ctx app.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var opts models.TargetQueryOptions
	if err := ctx.ShouldBindQuery(&opts); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	target, err := c.TargetService.Select(uint(id), opts)
	if err == errors.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, target)
}

func (c *TargetController) Delete(ctx app.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var opts models.TargetQueryOptions
	if err := ctx.ShouldBindQuery(&opts); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.TargetService.Delete(uint(id), opts)
	if err == errors.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
