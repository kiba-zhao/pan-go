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
