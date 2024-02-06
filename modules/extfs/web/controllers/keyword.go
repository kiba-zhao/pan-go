package controllers

import (
	"net/http"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/web"
)

type KeywordController struct {
	KeywordService *services.KeywordService
	DefaultLimit   int
}

// Init ...
func (c *KeywordController) Init(router web.Router) {
	r := router.Group("/extfs/keywords")
	r.GET("", c.Search)
}

// Search ...
func (c *KeywordController) Search(ctx web.Context) {
	var condition models.KeywordFindCondition
	var err error
	var keywords []models.Keyword
	if berr := ctx.ShouldBindQuery(&condition); berr != nil {
		keywords, err = c.KeywordService.Find(nil)
	} else {
		if condition.Limit == 0 {
			condition.Limit = c.DefaultLimit
		}
		keywords, err = c.KeywordService.Find(&condition)
	}

	if err == nil {
		ctx.JSON(http.StatusOK, keywords)
	} else {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

}
