// Define Remote Node controller for the web application
package remotenode

import (
	"net/http"
	"pan/app/web"
)

type RemoteNodeController struct {
	RemoteNodeService *RemoteNodeService
}

// SetupToWeb sets up the controller for the web application
//
// It sets up the following endpoints:
//
// - `GET /remote-nodes`: Search all remote nodes
//
// Example:
// type MyController struct {}
//
//	func (c *MyController) SetupToWeb(r WebRouter) error {
//		r.GET("/my-controller", c.MyMethod)
//	 ...
//		return nil
//	}
func (s *RemoteNodeController) SetupToWeb(router web.WebRouter) error {
	router.GET("/remote-nodes", s.Search)
	return nil
}

// Search returns a list of all remote nodes.
//
// It responds with a JSON array of RemoteNode objects and a total count
// in the X-Total-Count header.
//
// Example response:
// [
//
//	{
//		"id": <string>,
//		"name": <string>,
//		"description": <string>
//	},
//	{
//		"id": <string>,
//		"name": <string>,
//		"description": <string>
//	}
//
// ]
func (s *RemoteNodeController) Search(ctx web.WebContext) {

	total, nodes, err := s.RemoteNodeService.SelectAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, nodes)
}
