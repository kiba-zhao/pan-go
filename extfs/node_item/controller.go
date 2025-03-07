// Define node item controller for the web application
package nodeitem

import (
	"net/http"
	"pan/app/web"
	"strconv"
)

type NodeItemController struct {
	NodeItemService *NodeItemService
}

// SetupToWeb sets up the controller for the web application
//
// It sets up the following endpoints:
//
// - `GET /node-items`: Search all node items
// - `POST /node-items`: Create a new node item
// - `PATCH /node-items/:id`: Update a node item
// - `GET /node-items/:id`: Select a node item by its id
// - `DELETE /node-items/:id`: Delete a node item by its id
//
// Example:
// type MyController struct {}
//
//	func (c *MyController) SetupToWeb(r WebRouter) error {
//		r.GET("/my-controller", c.MyMethod)
//	 ...
//		return nil
//	}
func (c *NodeItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/node-items", c.Search)
	router.POST("/node-items", c.Create)
	router.PATCH("/node-items/:id", c.Update)
	router.GET("/node-items/:id", c.Select)
	router.DELETE("/node-items/:id", c.Delete)
	return nil
}

// Search returns a list of all node items.
//
// It responds with a JSON array of NodeItem objects and a total count
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
func (c *NodeItemController) Search(ctx web.WebContext) {

	total, items, err := c.NodeItemService.SelectAll()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

// Create creates a new node item.
//
// It expects a JSON object with the following fields in the request body:
//
// - `name`: The name of the node item.
// - `description`: The description of the node item.
//
// It responds with a JSON object with the newly created node item and a
// status code of 201 Created.
func (c *NodeItemController) Create(ctx web.WebContext) {

	var fields NodeItemFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	nodeItem, err := c.NodeItemService.Create(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, nodeItem)
}

// Update updates the node item with the specified id based on the given fields.
//
// It expects a JSON object with the following fields in the request body:
//
// - `name`: The name of the node item.
// - `description`: The description of the node item.
//
// It responds with a JSON object with the updated node item and a status code
// of 200 OK. If the node item was not found, it responds with a status code of
// 404 Not Found.
func (c *NodeItemController) Update(ctx web.WebContext) {

	var fields NodeItemFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nodeItem, err := c.NodeItemService.Update(fields, id)
	if err == ErrNodeItemNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

// Select returns the node item with the specified id.
//
// It responds with a JSON object with the node item and a status code of 200
// OK. If the node item was not found, it responds with a status code of 404 Not
// Found.
func (c *NodeItemController) Select(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	nodeItem, err := c.NodeItemService.Select(id)
	if c.NodeItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, nodeItem)
}

// Delete deletes the node item with the specified id.
//
// It responds with a status code of 204 No Content. If the node item was not
// found, it responds with a status code of 404 Not Found.
func (c *NodeItemController) Delete(ctx web.WebContext) {

	id, err := ExtractIdWithParam("id", ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = c.NodeItemService.Delete(id)
	if err == ErrNodeItemNotFound {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// ExtractIdWithParam extracts a uint id from the given web context parameter.
//
// It expects a parameter with the given name in the web context and attempts to
// parse it as a uint. If the parameter is not found or cannot be parsed, it
// returns an error.
//
// The extracted id is returned as a uint.
func ExtractIdWithParam(name string, ctx web.WebContext) (uint, error) {
	paramId := ctx.Param(name)
	id, err := strconv.ParseUint(paramId, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), err
}
