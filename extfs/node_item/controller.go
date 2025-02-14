package nodeitem

import (
	"net/http"
	"pan/app/web"
	"strconv"
)

type NodeItemController struct {
	NodeItemService *NodeItemService
}

func (c *NodeItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items", c.Search)
	router.POST("/node-items", c.Create)
	router.PATCH("/node-items/:id", c.Update)
	router.GET("/node-items/:id", c.Select)
	router.DELETE("/node-items/:id", c.Delete)
	return nil
}

func (c *NodeItemController) Search(ctx web.WebContext) {

	total, items, err := c.NodeItemService.SelectAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *NodeItemController) Create(ctx web.WebContext) {

	var fields NodeItemFields
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

func (c *NodeItemController) Update(ctx web.WebContext) {

	var fields NodeItemFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nodeItem, err := c.NodeItemService.Update(fields, id)
	if err == ErrNodeItemNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

func (c *NodeItemController) Select(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nodeItem, err := c.NodeItemService.Select(id)
	if c.NodeItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

func (c *NodeItemController) Delete(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.NodeItemService.Delete(id)
	if err == ErrNodeItemNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func ExtractIdWithParam(name string, ctx web.WebContext) (uint, error) {
	paramId := ctx.Param(name)
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), err
}
