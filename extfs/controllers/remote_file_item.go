package controllers

import (
	"bytes"
	"io"
	"net/http"
	appConstant "pan/app/constant"
	"pan/app/net"
	appNode "pan/app/node"
	"pan/extfs/models"
	"pan/extfs/services"

	"google.golang.org/protobuf/proto"
)

type RemoteFileItemController struct {
	RemoteFileItemService *services.RemoteFileItemService
}

func (c *RemoteFileItemController) SetupToWeb(router net.WebRouter) error {
	router.GET("/remote-file-items", c.Search)
	return nil
}

func (s *RemoteFileItemController) SetupToNode(router appNode.NodeRouter) error {
	router.Handle(services.RequestAllRemoteFileItems, s.SearchForNode)
	return nil
}

func (c *RemoteFileItemController) Search(ctx net.WebContext) {

	var condition models.RemoteFileItemSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := c.RemoteFileItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	net.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

func (c *RemoteFileItemController) SearchForNode(ctx *appNode.Context, next appNode.Next) error {

	req := ctx.Request()
	body, err := io.ReadAll(req.Body())
	if err != nil {
		ctx.ThrowError(appConstant.CodeBadRequest, err)
		return err
	}

	var condition models.RemoteFileItemRecordSearchCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(appConstant.CodeBadRequest, err)
		return err
	}

	fileItemList, err := c.RemoteFileItemService.SearchForNode(&condition)
	if err != nil {
		ctx.ThrowError(appConstant.CodeInternalError, err)
		return err
	}

	buffer, err := proto.Marshal(fileItemList)
	if err != nil {
		ctx.ThrowError(appConstant.CodeInternalError, err)
		return err
	}

	ctx.Respond(bytes.NewReader(buffer))
	return err
}
