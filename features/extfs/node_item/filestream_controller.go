// Define node file stream controller for the web application
package nodeitem

import (
	"errors"
	"net/http"
	"pan/lib/web"
)

type NodeFileStreamController struct {
	NodeFilePathService *NodeFilePathService
}

// SetupToWeb sets up the routes for the NodeFileStreamController in the web application.
//
// It configures the following endpoints:
//
// - `GET /node-items/:id/_stream/*filepath`: Triggers the Select method to retrieve a file by its filepath associated with a node item by its ID.
func (c *NodeFileStreamController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items/:id/_stream/*filepath", c.Select)
	return nil
}

// Select retrieves a file associated with a node item by its ID and filepath.
//
// It extracts the node item ID from the URL parameters and the filepath from the URL parameter 'filepath'.
// The method interacts with the NodeFilePathService to fetch the real file path without folder information.
// If the ID is invalid or the file path does not exist, it returns a 400 Bad Request or 404 Not Found error, respectively.
// If the file path cannot be accessed without folder information, it returns a 403 Forbidden error.
// On successful retrieval, it serves the file to the client.
func (ctrl *NodeFileStreamController) Select(ctx web.WebContext) {
	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	filePath := ctx.Param("filepath")

	realFilePath, err := ctrl.NodeFilePathService.SelectWithoutFolder(id, filePath[1:])
	if ctrl.NodeFilePathService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if errors.Is(err, ErrNodeFilePathWithoutFolder) {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.File(realFilePath)
}
