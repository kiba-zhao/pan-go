package controllers

import (
	"net/http"

	"pan/services"
	"pan/web"
)

type PeerStatController struct {
	PeerStat *services.PeerStatService
}

// Route ...
func (c *PeerStatController) Init(router web.Router) {
	r := router.Group("/base/peer-stat")
	r.GET("/:id", c.Get)
}

// Get ...
func (c *PeerStatController) Get(ctx web.Context) {
	id := ctx.Param("id")
	stat, err := c.PeerStat.FindOne(id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.JSON(http.StatusOK, stat)
}
