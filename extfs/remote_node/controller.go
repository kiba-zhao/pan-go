package remotenode

import (
	"net/http"
	"pan/app/web"
)

type RemoteNodeController struct {
	RemoteNodeService *RemoteNodeService
}

func (s *RemoteNodeController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote-nodes", s.Search)
	return nil
}

func (s *RemoteNodeController) Search(ctx web.WebContext) {

	total, nodes, err := s.RemoteNodeService.SelectAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, nodes)
}
