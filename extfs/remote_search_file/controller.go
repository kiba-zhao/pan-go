// Define controller for web application
//
// It handles web api request for remote search file
package remotesearchfile

import (
	"errors"
	"net/http"
	"pan/app/web"
)

var ErrRemoteSearchFileInvalidPeerID = errors.New("remotesearchfile.RemoteSearchFileController Error: Invalid Peer ID")

type RemoteSearchFileController struct {
	RemoteSearchFileService *RemoteSearchFileService
}

// SetupToWeb sets up the controller for the web application
//
// It sets up the following endpoint:
//
// - `GET /remotes/:peerId/search-files`: Search all remote search files with the given peer ID
func (c *RemoteSearchFileController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remotes/:peerId/search-files", c.Search)
	return nil
}

// Search handles the HTTP request for searching remote search files associated with a specific peer ID.
//
// It extracts the peer ID from the URL parameters and binds the search condition from the request body.
// If the peer ID is not provided, or if binding the search condition fails, it aborts the request with a
// 400 Bad Request status. If an "If-Match" header is present and the search condition hash is empty, it
// uses the header value as the hash.
//
// The method then calls the RemoteSearchFileService to perform the search using the peer ID and search
// condition. If the search fails, it returns a 500 Internal Server Error status.
//
// On success, it sets the X-Total-Count and ETag headers with the total count and ETag value, respectively,
// and responds with a 200 OK status along with the list of remote search files in JSON format.

func (c *RemoteSearchFileController) Search(ctx web.WebContext) {
	paramPeerId := ctx.Param("peerId")
	if len(paramPeerId) <= 0 {
		ctx.AbortWithError(http.StatusBadRequest, ErrRemoteSearchFileInvalidPeerID)
		return
	}

	var condition RemoteSearchFileSearchCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ifMatch := ctx.Request.Header.Get("If-Match")
	if len(condition.Hash) <= 0 && len(ifMatch) > 0 {
		condition.Hash = ifMatch
	}

	total, items, etag, err := c.RemoteSearchFileService.Search(paramPeerId, condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	web.SetCountHeaderForWeb(ctx, total)
	web.SetETagHeaderForWeb(ctx, etag)
	ctx.JSON(http.StatusOK, items)
}
