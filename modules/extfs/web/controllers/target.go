package controllers

import (
	"pan/modules/extfs/services"
	"pan/web"
)

type TargetController struct {
	TargetService *services.TargetService
}

// Init
func (c *TargetController) Init(router web.Router) {
	r := router.Group("/extfs/targets")
	r.GET("", c.Search)
}

// Search
func (c *TargetController) Search(ctx web.Context) {

}
