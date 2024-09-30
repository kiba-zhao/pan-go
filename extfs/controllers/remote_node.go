package controllers

import (
	"net/http"
	"pan/app/net"
	"pan/extfs/services"
)

type RemoteNodeController struct {
	RemoteNodeService *services.RemoteNodeService
}

func (s *RemoteNodeController) SetupToWeb(router net.WebRouter) error {
	router.GET("/remote-nodes", s.Search)
	return nil
}

func (s *RemoteNodeController) Search(ctx net.WebContext) {

	total, nodes, err := s.RemoteNodeService.SelectAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, nodes)
}
