package controllers

import "pan/app/net"

type NodeController struct {
}

func (c *NodeController) SetupToWeb(router net.WebRouter) error {
	router.GET("/nodes", c.Search)
	return nil
}

func (c *NodeController) Search(ctx net.WebContext) {
}
