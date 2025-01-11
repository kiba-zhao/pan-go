package nodesearchfile

import (
	"net/http"
	"pan/app/web"
)

type NodeSearchFileController struct {
	NodeSearchFileService *NodeSearchFileService
}

func (c *NodeSearchFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/search-files", c.Search)
	return nil
}

func (c *NodeSearchFileController) Search(ctx web.WebContext) {

	var condition NodeSearchFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ifMatch := ctx.Request.Header.Get("If-Match")
	if len(condition.Hash) <= 0 && len(ifMatch) > 0 {
		condition.Hash = ifMatch
	}

	total, items, etag, err := c.NodeSearchFileService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	web.SetCountHeaderForWeb(ctx, total)
	web.SetETagHeaderForWeb(ctx, etag)
	ctx.JSON(http.StatusOK, items)
}
