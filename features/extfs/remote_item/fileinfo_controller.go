package remoteitem

import (
	"net/http"
	nodeitem "pan/features/extfs/node_item"
	"pan/lib/web"
)

type RemoteFileInfoController struct {
	RemoteFileInfoService *RemoteFileInfoService
}

func (ctrl *RemoteFileInfoController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remotes/:peerId/remote-items/:id/_files", ctrl.Search)
	router.GET("/remotes/:peerId/remote-items/:id/_files/*filepath", ctrl.Select)
	return nil
}

func (ctrl *RemoteFileInfoController) Search(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := nodeitem.ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	parentPath := ctx.Query("parentPath")

	total, infos, err := ctrl.RemoteFileInfoService.Search(peerId, id, parentPath)
	if ctrl.RemoteFileInfoService.IsNotExist(err) {
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

func (ctrl *RemoteFileInfoController) Select(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := nodeitem.ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	filePath := ctx.Param("filepath")

	remoteFileInfo, err := ctrl.RemoteFileInfoService.Select(peerId, id, filePath[1:])
	if ctrl.RemoteFileInfoService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, remoteFileInfo)
}
