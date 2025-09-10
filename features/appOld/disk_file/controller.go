package diskfile

import (
	"net/http"
	"pan/lib/web"
)

type DiskFileController struct {
	DiskFileService *DiskFileService
}

func (c *DiskFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/disk-files", c.Search)
	return nil
}

func (c *DiskFileController) Search(ctx web.WebContext) {
	var conditions DiskFileSearchCondition
	if err := ctx.ShouldBind(&conditions); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := c.DiskFileService.Search(conditions)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err == ErrDiskFileParentPathConflict {
		ctx.AbortWithError(http.StatusConflict, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	if total == 0 {
		items = []DiskFile{}
	}
	ctx.JSON(http.StatusOK, items)
}
