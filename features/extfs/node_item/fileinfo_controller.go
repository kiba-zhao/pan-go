// Define node file info controller for the web application
package nodeitem

import (
	"net/http"
	"pan/lib/web"
)

type NodeFileInfoController struct {
	NodeFileInfoService *NodeFileInfoService
}

// SetupToWeb sets up the routes for the NodeFileInfoController in the web application.
//
// It configures the following endpoints:
//
// - `GET /node-items/:id/_files`: Triggers the Search method to list files associated with a specific node item by its ID.
// - `GET /node-items/:id/_files/*filepath`: Triggers the Select method to retrieve a specific file by its filepath.
//

func (c *NodeFileInfoController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items/:id/_files", c.Search)
	router.GET("/node-items/:id/_files/*filepath", c.Select)
	return nil
}

// Search retrieves a list of file information associated with a specific node item by its ID.
//
// It extracts the node item ID from the URL parameters and an optional parentPath
// from the query parameters. The method interacts with the NodeFileInfoService to
// fetch the file information. If the ID is invalid or the item does not exist, it
// returns a 400 Bad Request or 404 Not Found error, respectively. On successful
// retrieval, it sets the X-Total-Count header with the total number of items and
// responds with a JSON array of file information.
//
// Parameters:
// - ctx: The web context containing request and response information.

func (c *NodeFileInfoController) Search(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	parentPath := ctx.Query("parentPath")

	total, infos, err := c.NodeFileInfoService.Search(id, parentPath)
	if c.NodeFileInfoService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, infos)
}

// Select retrieves a specific file information associated with a specific node item by its ID and filepath.
//
// It extracts the node item ID from the URL parameters and the filepath from the URL parameter 'filepath'.
// The method interacts with the NodeFileInfoService to fetch the file information. If the ID is invalid or the item does not exist, it
// returns a 400 Bad Request or 404 Not Found error, respectively. On successful retrieval, it responds with a JSON object of file information.
func (c *NodeFileInfoController) Select(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filePath := ctx.Param("filepath")

	nodeFileInfo, err := c.NodeFileInfoService.Select(id, filePath[1:])
	if err == ErrNodeFileInfoInvalidID {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.NodeFileInfoService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeFileInfo)
}
