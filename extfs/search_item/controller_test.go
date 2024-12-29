package searchitem_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"pan/app/web"
	"strconv"
	"testing"

	searchitem "pan/extfs/search_item"
	mocked "pan/mocks/pan/extfs/search_item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchItemController(t *testing.T) {
	setup := func() (web.WebApp, *searchitem.SearchItemController) {
		app := web.NewWebApp()
		ctrl := &searchitem.SearchItemController{}
		ctrl.SetupToWeb(app)

		ctrl.SearchItemService = &searchitem.SearchItemService{}
		return app, ctrl
	}

	t.Run("GET /search-items?q=", func(t *testing.T) {
		app, ctrl := setup()

		searchItemRepo := &mocked.MockSearchItemRepository{}
		defer searchItemRepo.AssertExpectations(t)
		ctrl.SearchItemService.SearchItemRepo = searchItemRepo

		query := "query"
		total := int64(10)
		items := []searchitem.SearchItem{
			{ID: 1, Query: "query 1"},
			{ID: 2, Query: "query 2"},
			{ID: 3, Query: "query 3"},
		}
		searchItemRepo.On("Search", mock.Anything).Once().Return(total, items, nil).Run(func(args mock.Arguments) {
			condition := args.Get(0).(searchitem.SearchItemCondition)
			assert.Equal(t, query, condition.Query)
			assert.Equal(t, int(total), condition.RangeEnd)
		})

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/search-items", nil)
		q := r.URL.Query()
		q.Add("q", query)
		q.Add("_end", strconv.FormatInt(total, 10))
		r.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))

		var results []searchitem.SearchItem
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.NoError(t, err)
		assert.Equal(t, items, results)
	})

	t.Run("DELETE /search-items/:id", func(t *testing.T) {
		app, ctrl := setup()

		searchItemRepo := &mocked.MockSearchItemRepository{}
		defer searchItemRepo.AssertExpectations(t)
		ctrl.SearchItemService.SearchItemRepo = searchItemRepo

		id := uint64(1)
		searchItemRepo.On("Delete", id).Once().Return(nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/search-items/%d", id)
		r := httptest.NewRequest("DELETE", url, nil)
		app.ServeHTTP(w, r)

		assert.Equal(t, 204, w.Code)

	})

	t.Run("POST /search-items", func(t *testing.T) {
		app, ctrl := setup()

		searchItemRepo := &mocked.MockSearchItemRepository{}
		defer searchItemRepo.AssertExpectations(t)
		ctrl.SearchItemService.SearchItemRepo = searchItemRepo

		query := "query"
		model := searchitem.SearchItem{Query: query}
		newModel := model
		newModel.ID = 1

		searchItemRepo.On("SelectOrCreate", model).Once().Return(newModel, nil)

		fields := searchitem.SearchItemFields{Query: query}
		jsonData, _ := json.Marshal(fields)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/search-items", bytes.NewReader(jsonData))
		r.Header.Set("Content-Type", "application/json")
		app.ServeHTTP(w, r)

		assert.Equal(t, 201, w.Code)

		var result searchitem.SearchItem
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, newModel, result)
	})
}
