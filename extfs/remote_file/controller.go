package remotefile

import (
	"net/http"
	"pan/app/web"
)

type RemoteFileController struct {
	RemoteFileService *RemoteFileService
}

func (c *RemoteFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote-files", c.Search)
	return nil
}

func (c *RemoteFileController) Search(ctx web.WebContext) {

	var condition RemoteFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.RemoteFileService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}
