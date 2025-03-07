// Define remote item controller for the web application
package remoteitem

import (
	"errors"
	"net/http"
	"pan/app/web"
	nodeitem "pan/extfs/node_item"
)

type RemoteItemController struct {
	RemoteItemService *RemoteItemService
}

// SetupToWeb sets up the controller for the web application
//
// It sets up the following endpoints:
//
// - `GET /remotes/:peerId/remote-items`: Search all remote items with the given peer ID
// - `GET /remotes/:peerId/remote-items/:id`: Select a remote item by its ID and peer ID
func (ctrl *RemoteItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remotes/:peerId/remote-items", ctrl.Search)
	router.GET("/remotes/:peerId/remote-items/:id", ctrl.Select)
	return nil
}

// Search handles the HTTP request to retrieve all remote items associated with a peer ID.
//
// It extracts the peer ID from the URL parameters and queries the RemoteItemService
// for the list of items. If the extraction or service query fails, it returns an appropriate
// HTTP error response. On success, it sets the total count in the response header and
// returns the list of items as a JSON array with a 200 OK status.

func (ctrl *RemoteItemController) Search(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	total, items, err := ctrl.RemoteItemService.Search(peerId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

// Select retrieves a specific remote item by its ID and peer ID.
//
// It extracts the peer ID and item ID from the URL parameters and interacts
// with the RemoteItemService to fetch the item. If either the peer ID or
// item ID is invalid, or if the item does not exist, it returns a 400
// Bad Request or 404 Not Found error, respectively. On successful retrieval,
// it responds with a JSON object of the remote item.

func (c *RemoteItemController) Select(ctx web.WebContext) {
	peerId, err := ExtractPeerIdWithParam("peerId", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	id, err := nodeitem.ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	remoteItem, err := c.RemoteItemService.Select(peerId, id)
	if errors.Is(err, ErrRemoteItemInvalidID) {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if c.RemoteItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, remoteItem)
}

// ExtractPeerIdWithParam extracts a string peer ID from the given web context parameter.
//
// It expects a parameter with the given name in the web context and attempts to
// parse it as a string. If the parameter is not found or is empty, it returns an
// error.
//
// The extracted peer ID is returned as a string.
func ExtractPeerIdWithParam(name string, ctx web.WebContext) (string, error) {
	peerId := ctx.Param("peerId")
	if len(peerId) <= 0 {
		return "", errors.New("invalid peer id")
	}
	return peerId, nil
}
