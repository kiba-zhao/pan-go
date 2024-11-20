package appnode

import (
	"net/http"

	"pan/app/web"
	"strconv"
)

type AppNodeController struct {
	AppNodeService *AppNodeService
}

func (c *AppNodeController) SetupToWeb(router web.WebRouter) error {
	router.GET("/nodes", c.Search)
	router.GET("/nodes/:id", c.Select)
	router.DELETE("/nodes/:id", c.Delete)
	router.POST("/nodes", c.Create)
	router.PATCH("/nodes/:id", c.Update)
	return nil
}

func (c *AppNodeController) Search(ctx web.WebContext) {
	var conditions AppNodeSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.AppNodeService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *AppNodeController) Select(ctx web.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	node, err := c.AppNodeService.Select(uint(id))
	if err == ErrAppNodeNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, node)
}

func (c *AppNodeController) Delete(ctx web.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.AppNodeService.Delete(uint(id))
	if err == ErrAppNodeNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *AppNodeController) Create(ctx web.WebContext) {
	var fields AppNodeFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	node, err := c.AppNodeService.Create(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, node)
}

func (c *AppNodeController) Update(ctx web.WebContext) {
	var fields AppNodeFields
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
	node, err := c.AppNodeService.Update(uint(id), fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, node)
}
