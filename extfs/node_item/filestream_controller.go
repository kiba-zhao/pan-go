package nodeitem

import (
	"net/http"
	"pan/app/web"
)

type NodeFileStreamController struct {
	NodeFilePathService *NodeFilePathService
}

func (c *NodeFileStreamController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items/:id/stream/*filepath", c.Select)
	return nil
}

func (ctrl *NodeFileStreamController) Select(ctx web.WebContext) {
	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filePath := ctx.Param("filepath")

	realFilePath, err := ctrl.NodeFilePathService.Select(id, filePath[1:])
	if ctrl.NodeFilePathService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.File(realFilePath)
}
