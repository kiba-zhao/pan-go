package nodefile

import (
	"net/http"
	"pan/app/web"
)

type NodeFileController struct {
	NodeFileService *NodeFileService
}

func (c *NodeFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-files", c.Search)
	router.GET("/node-files/:id", c.Select)
	return nil
}

func (c *NodeFileController) Search(ctx web.WebContext) {

	var condition NodeFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.NodeFileService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *NodeFileController) Select(ctx web.WebContext) {
	paramId := ctx.Param("id")
	nodeFile, err := c.NodeFileService.Select(paramId)
	if err == ErrNodeFileInvalidID {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.NodeFileService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeFile)
}
