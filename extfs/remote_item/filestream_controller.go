package remoteitem

import (
	"net/http"
	"pan/app/web"
	nodeitem "pan/extfs/node_item"
)

type RemoteFileStreamController struct {
	RemoteFileStreamService *RemoteFileStreamService
}

func (ctrl *RemoteFileStreamController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remotes/:peerId/remote-items/:id/_stream/*filepath", ctrl.Select)
	return nil
}

func (ctrl *RemoteFileStreamController) Select(ctx web.WebContext) {
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
	reader, err := ctrl.RemoteFileStreamService.Read(peerId, id, filePath[1:])
	if ctrl.RemoteFileStreamService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	defer reader.Close()
	http.ServeContent(ctx.Writer, ctx.Request, reader.Name(), reader.ModTime(), reader)
}
