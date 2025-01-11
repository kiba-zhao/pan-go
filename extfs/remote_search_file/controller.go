package remotesearchfile

import (
	"errors"
	"net/http"
	"pan/app/web"
)

var ErrRemoteSearchFileInvalidPeerID = errors.New("remotesearchfile.RemoteSearchFileController Error: Invalid Peer ID")

type RemoteSearchFileController struct {
	RemoteSearchFileService *RemoteSearchFileService
}

func (c *RemoteSearchFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote/:peerId/search-files", c.Search)
	return nil
}

func (c *RemoteSearchFileController) Search(ctx web.WebContext) {
	paramPeerId := ctx.Param("peerId")
	if len(paramPeerId) <= 0 {
		ctx.AbortWithError(http.StatusBadRequest, ErrRemoteSearchFileInvalidPeerID)
		return
	}

	var condition RemoteSearchFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ifMatch := ctx.Request.Header.Get("If-Match")
	if len(condition.Hash) <= 0 && len(ifMatch) > 0 {
		condition.Hash = ifMatch
	}

	total, items, etag, err := c.RemoteSearchFileService.Search(paramPeerId, condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	web.SetCountHeaderForWeb(ctx, total)
	web.SetETagHeaderForWeb(ctx, etag)
	ctx.JSON(http.StatusOK, items)
}
