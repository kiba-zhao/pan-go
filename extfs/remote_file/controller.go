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
	router.GET("/remote-files/:id", c.Select)
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

func (c *RemoteFileController) Select(ctx web.WebContext) {
	paramId := ctx.Param("id")
	remoteFile, err := c.RemoteFileService.Select(paramId)
	if err == ErrRemoteFileInvalidID {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.RemoteFileService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, remoteFile)
}
