package nodeitem

import (
	"net/http"
	"pan/app/web"
)

type NodeFileInfoController struct {
	NodeFileInfoService *NodeFileInfoService
}

func (c *NodeFileInfoController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items/:id/files/*filepath", c.Search)
	router.GET("/node-items/:id/file/*filepath", c.Select)
	return nil
}

func (c *NodeFileInfoController) Search(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filePath := ctx.Param("filepath")

	total, infos, err := c.NodeFileInfoService.Search(id, filePath[1:])
	if c.NodeFileInfoService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, infos)
}

func (c *NodeFileInfoController) Select(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filePath := ctx.Param("filepath")

	nodeFileInfo, err := c.NodeFileInfoService.Select(id, filePath[1:])
	if err == ErrNodeFileInfoInvalidID {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.NodeFileInfoService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeFileInfo)
}
