package controllers

import (
	"net/http"
	"pan/app/constant"
	"pan/app/net"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
)

type NodeItemController struct {
	NodeItemService *services.NodeItemService
}

func (c *NodeItemController) SetupToWeb(router net.WebRouter) error {
	router.POST("/node-items", c.Create)
	router.PATCH("/node-items/:id", c.Update)
	router.GET("/node-items/:id", c.Select)
	router.DELETE("/node-items/:id", c.Delete)
	return nil
}

func (c *NodeItemController) Create(ctx net.WebContext) {

	var fields models.NodeItemFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	nodeItem, err := c.NodeItemService.Create(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, nodeItem)
}

func (c *NodeItemController) Update(ctx net.WebContext) {

	var fields models.NodeItemFields
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

	nodeItem, err := c.NodeItemService.Update(fields, uint(id))
	if err == constant.ErrConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

func (c *NodeItemController) Select(ctx net.WebContext) {

	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nodeItem, err := c.NodeItemService.Select(uint(id))
	if err == constant.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

func (c *NodeItemController) Delete(ctx net.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.NodeItemService.Delete(uint(id))
	if err == constant.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
