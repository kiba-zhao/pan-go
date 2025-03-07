package searchitem

import (
	"net/http"
	"pan/app/web"
	"strconv"
)

type SearchItemController struct {
	SearchItemService *SearchItemService
}

// SetupToWeb sets up the routes for the SearchItemController in the web application.
//
// It configures the following endpoints:
//
// - `GET /search-items`: Triggers the Search method to retrieve search items based on conditions.
// - `DELETE /search-items/:id`: Triggers the Delete method to remove a search item by its ID.
// - `POST /search-items`: Triggers the Create method to add a new search item.

func (ctrl *SearchItemController) SetupToWeb(router web.WebRouter) error {
	router.GET("/search-items", ctrl.Search)
	router.DELETE("/search-items/:id", ctrl.Delete)
	router.POST("/search-items", ctrl.Create)
	return nil
}

// Search retrieves search items based on conditions.
//
// It expects a JSON object with the following fields in the request body:
//
// - `q`: The search query.
// - `query`: The name of the query.
// - `_end`: The end of the range of search item IDs to retrieve.
//
// It responds with a JSON array of SearchItem objects, a total count
// in the X-Total-Count header, and a status code of 200 OK. If the
// search item service returns an error, it responds with the error and
// a status code of 500 Internal Server Error.
func (ctrl *SearchItemController) Search(ctx web.WebContext) {
	var condition SearchItemCondition
	if err := ctx.ShouldBind(&condition); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	total, items, err := ctrl.SearchItemService.Search(condition)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	web.SetCountHeaderForWeb(ctx, total)
	ctx.JSON(http.StatusOK, items)
}

// Delete removes a search item by its ID.
//
// It extracts the search item ID from the URL parameters and attempts to delete the corresponding search item.
// If the ID is invalid, it returns a 400 Bad Request error.
// If the search item does not exist, it returns a 404 Not Found error.
// If an internal error occurs during deletion, it returns a 500 Internal Server Error.
// On successful deletion, it responds with a status code of 204 No Content.

func (ctrl *SearchItemController) Delete(ctx web.WebContext) {
	paramId := ctx.Param("id")
	id, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = ctrl.SearchItemService.Delete(id)
	if ctrl.SearchItemService.IsNotExist(err) {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Create adds a new search item if it does not exist or selects an existing search item otherwise.
//
// It expects a JSON object with the following fields in the request body:
//
// - `q`: The search query.
// - `query`: The name of the query.
//
// It responds with a JSON object of the newly created search item and a status code of 201 Created. If the
// search item service returns an error, it responds with the error and a status code of 500 Internal Server Error.
func (ctrl *SearchItemController) Create(ctx web.WebContext) {
	var fields SearchItemFields
	if err := ctx.ShouldBind(&fields); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	searchItem, err := ctrl.SearchItemService.SelectOrCreate(fields)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, searchItem)
}
