package controllers

import (
	"net/http"
	"pan/app/constant"
	"pan/app/models"
	"pan/app/net"
	"pan/app/services"
	"strconv"
)

type NodeController struct {
	NodeService *services.NodeService
}

func (c *NodeController) SetupToWeb(router net.WebRouter) error {
	router.GET("/nodes", c.Search)
	router.GET("/nodes/:id", c.Select)
	router.DELETE("/nodes/:id", c.Delete)
	router.POST("/nodes", c.Create)
	router.PATCH("/nodes/:id", c.Update)
	return nil
}

func (c *NodeController) Search(ctx net.WebContext) {
	var conditions models.NodeSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.NodeService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *NodeController) Select(ctx net.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	node, err := c.NodeService.Select(uint(id))
	if err == constant.ErrNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, node)
}

func (c *NodeController) Delete(ctx net.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.NodeService.Delete(uint(id))
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

func (c *NodeController) Create(ctx net.WebContext) {
	var fields models.NodeFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	node, err := c.NodeService.Create(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, node)
}

func (c *NodeController) Update(ctx net.WebContext) {
	var fields models.NodeFields
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
	node, err := c.NodeService.Update(uint(id), fields)
	if err == constant.ErrConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, node)
}
