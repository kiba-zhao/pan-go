package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pan/modules/extfs/controllers"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/web"
	"strconv"
	"testing"

	mockedRepo "pan/mocks/pan/modules/extfs/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTarget(t *testing.T) {

	setup := func() *controllers.TargetController {
		ctrl := new(controllers.TargetController)
		ctrl.TargetService = new(services.TargetService)
		ctrl.TargetService.TargetFileService = new(services.TargetFileService)
		return ctrl
	}

	t.Run("Search", func(t *testing.T) {
		app := web.NewApp()
		ctrl := setup()
		ctrl.MountWithWeb(app)

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		var condition models.TargetSearchCondition
		condition.Keyword = "Keyword"
		condition.Limit = 112
		condition.Offset = 12

		total := int64(123)
		targets := []models.Target{
			models.Target{Name: "name 1"},
			models.Target{Name: "name 2"},
		}
		targetRepo.On("Search", mock.Anything).Once().Return(total, targets, nil).Run(func(args mock.Arguments) {
			c := args.Get(0).(*models.TargetSearchCondition)

			assert.Equal(t, condition.Keyword, c.Keyword)
			assert.Equal(t, condition.Limit, c.Limit)
			assert.Equal(t, condition.Offset, c.Offset)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)

		q := req.URL.Query()
		q.Add("keyword", condition.Keyword)
		q.Add("limit", strconv.Itoa(condition.Limit))
		q.Add("offset", strconv.Itoa(condition.Offset))
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var results models.TargetSearchResult
		bodyBytes := w.Body.Bytes()
		err := json.Unmarshal(bodyBytes, &results)
		assert.Nil(t, err)
		assert.Equal(t, total, results.Total)
		assert.Equal(t, targets, results.Items)
	})

	t.Run("Search Without condition", func(t *testing.T) {
		app := web.NewApp()
		ctrl := setup()
		ctrl.MountWithWeb(app)

		targetRepo := new(mockedRepo.MockTargetRepository)
		defer targetRepo.AssertExpectations(t)
		ctrl.TargetService.TargetRepo = targetRepo

		total := int64(123)
		targets := []models.Target{
			models.Target{Name: "name 1"},
			models.Target{Name: "name 2"},
		}
		targetRepo.On("Search", mock.Anything).Once().Return(total, targets, nil).Run(func(args mock.Arguments) {
			c := args.Get(0).(*models.TargetSearchCondition)

			assert.Equal(t, "", c.Keyword)
			assert.Equal(t, 100, c.Limit)
			assert.Equal(t, 0, c.Offset)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var results models.TargetSearchResult
		bodyBytes := w.Body.Bytes()
		err := json.Unmarshal(bodyBytes, &results)
		assert.Nil(t, err)
		assert.Equal(t, total, results.Total)
		assert.Equal(t, targets, results.Items)
	})

	t.Run("Search With Error", func(t *testing.T) {
		app := web.NewApp()
		ctrl := setup()
		ctrl.MountWithWeb(app)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/targets", nil)

		q := req.URL.Query()
		q.Add("limit", "limit")
		q.Add("offset", "offset")
		req.URL.RawQuery = q.Encode()
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
