package controllers

import "pan/web"

type FileController struct {
}

// Init ...
func (c *FileController) Init(router web.Router) {

	r := router.Group("/extfs/files")
	r.GET("", c.Search)
}

// Search ...
func (c *FileController) Search(ctx web.Context) {

}
