package controllers

import (
	"bytes"
	"net/http"
	"pan/app/constant"
	"pan/app/net"
	appNode "pan/app/node"
	"pan/extfs/models"
	"pan/extfs/services"

	"google.golang.org/protobuf/proto"
)

type RemoteNodeItemController struct {
	RemoteNodeItemService *services.RemoteNodeItemService
}

func (s *RemoteNodeItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/remote-node-items", s.Search)
	return nil
}

func (s *RemoteNodeItemController) SetupToNode(router appNode.NodeRouter) error {
	router.Handle(services.RequestAllRemoteItems, s.SearchForNode)
	return nil
}

func (s *RemoteNodeItemController) Search(ctx net.WebContext) {
	var condition models.RemoteNodeItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := s.RemoteNodeItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (s *RemoteNodeItemController) SearchForNode(ctx *appNode.Context, next appNode.Next) error {

	recordList, err := s.RemoteNodeItemService.SelectAllForNode()
	if err != nil {
		ctx.ThrowError(constant.CodeInternalError, err)
		return err
	}

	buffer, err := proto.Marshal(&recordList)
	if err != nil {
		ctx.ThrowError(constant.CodeInternalError, err)
		return err
	}

	ctx.Respond(bytes.NewReader(buffer))
	return err
}
