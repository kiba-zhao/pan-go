package nodesearchfile

import (
	"net/http"
	"pan/lib/web"
)

type NodeSearchFileController struct {
	NodeSearchFileService *NodeSearchFileService
}

// SetupToWeb sets up the controller for the web application
//
// It sets up the following endpoint:
//
// - `GET /search-files`: Search all files with the given condition
func (c *NodeSearchFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/search-files", c.Search)
	return nil
}

// Search handles the HTTP GET request to search for files based on a given set of conditions.
//
// It binds the request parameters to a NodeSearchFileSearchCondition structure. If a 'Hash' is not provided
// in the request, it checks for an 'If-Match' header to use as a fallback hash. The function then calls
// the NodeSearchFileService to perform the search, which returns the total count, a list of NodeSearchFile
// items, and an ETag.
//
// If the search fails because the resources do not exist, it responds with a 404 Not Found error.
// For any other errors, it responds with a 500 Internal Server Error. If the search is successful,
// it sets the X-Total-Count and ETag headers and responds with a 200 OK status and the list of found items
// in JSON format. If no items are found, the response contains an empty list.

func (c *NodeSearchFileController) Search(ctx web.WebContext) {

	var condition NodeSearchFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ifMatch := ctx.Request.Header.Get("If-Match")
	if len(condition.Hash) <= 0 && len(ifMatch) > 0 {
		condition.Hash = ifMatch
	}

	total, items, etag, err := c.NodeSearchFileService.Search(condition)

	if c.NodeSearchFileService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if items == nil {
		items = make([]NodeSearchFile, 0)
	}
	web.SetCountHeaderForWeb(ctx, total)
	web.SetETagHeaderForWeb(ctx, etag)
	ctx.JSON(http.StatusOK, items)
}
