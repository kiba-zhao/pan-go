package controllers

import (
	"net/http"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/web"
)

type TagController struct {
	TagService   *services.TagService
	DefaultLimit int
}

// Init ...
func (c *TagController) Init(router web.Router) {
	r := router.Group("/extfs/tags")
	r.GET("", c.Search)
}

// Get ...
func (c *TagController) Search(ctx web.Context) {
	var condition models.TagFindCondition
	var err error
	var tags []models.Tag
	if berr := ctx.ShouldBindQuery(&condition); berr != nil {
		tags, err = c.TagService.Find(nil)
	} else {
		if condition.Limit == 0 {
			condition.Limit = c.DefaultLimit
		}
		tags, err = c.TagService.Find(&condition)
	}

	if err == nil {
		ctx.JSON(http.StatusOK, tags)
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

}
