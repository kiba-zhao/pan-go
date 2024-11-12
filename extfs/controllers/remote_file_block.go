package controllers

import (
	"io"
	appConstant "pan/app/constant"
	appNode "pan/app/node"
	"pan/extfs/models"
	"pan/extfs/services"

	"google.golang.org/protobuf/proto"
)

type RemoteFileBlockController struct {
	RemoteFileBlockService *services.RemoteFileBlockService
}

func (c *RemoteFileBlockController) SetupToNode(router appNode.NodeRouter) error {
	router.Handle(services.RequestRemoteFileBlock, c.SelectForNode)
	return nil
}

func (c *RemoteFileBlockController) SelectForNode(ctx *appNode.Context, next appNode.Next) error {
	req := ctx.Request()
	body, err := io.ReadAll(req)
	if err != nil {
		ctx.ThrowError(appConstant.CodeBadRequest, err)
		return err
	}

	var condition models.RemoteFileBlockSelectCondition
	err = proto.Unmarshal(body, &condition)
	if err != nil {
		ctx.ThrowError(appConstant.CodeBadRequest, err)
		return err
	}

	res, err := c.RemoteFileBlockService.SelectForNode(&condition)
	if err != nil {
		ctx.ThrowError(appConstant.CodeInternalError, err)
		return err
	}

	ctx.Respond(res)
	return err
}
