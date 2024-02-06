package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pan/mocks/pan/modules/extfs/repositories"
	"pan/modules/extfs/models"
	"pan/modules/extfs/services"
	"pan/modules/extfs/web/controllers"
	"pan/web"
	"strconv"
	"testing"

	mocked "pan/mocks/pan/modules/extfs/repositories"

	"github.com/stretchr/testify/assert"
)

// TestTags ...
func TestTags(t *testing.T) {
	setup := func(app *web.App, repo *repositories.MockKeywordRepository, limit int) {
		ctrl := new(controllers.KeywordController)
		ctrl.KeywordService = new(services.KeywordService)
		ctrl.KeywordService.KeywordRepo = repo
		ctrl.DefaultLimit = limit

		ctrl.Init(app)
	}

	t.Run("GET /extfs/keywords", func(t *testing.T) {
		limit := 101
		keywordRepo := new(mocked.MockKeywordRepository)
		app := web.NewApp()
		setup(app, keywordRepo, 101)

		keyword := new(models.Keyword)
		keyword.Name = "name"
		keywords := []models.Keyword{*keyword}

		condition := new(models.KeywordFindCondition)
		condition.Limit = limit
		keywordRepo.On("Find", condition).Once().Return(keywords, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/extfs/keywords", nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body, err := json.Marshal(keywords)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), w.Body.String())

		keywordRepo.AssertExpectations(t)
	})

	t.Run("GET /extfs/keywords with special name", func(t *testing.T) {

		limit := 100
		keywordRepo := new(mocked.MockKeywordRepository)
		app := web.NewApp()
		setup(app, keywordRepo, limit)

		qname := "ta"
		qlimit := 10
		qoffset := 0
		keyword := new(models.Keyword)
		keyword.Name = "keyword name"
		keywords := []models.Keyword{*keyword}

		condition := new(models.KeywordFindCondition)
		condition.Keyword = qname
		condition.Limit = qlimit
		condition.Offset = qoffset
		keywordRepo.On("Find", condition).Once().Return(keywords, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/extfs/keywords", nil)
		q := req.URL.Query()
		q.Add("keyword", qname)
		q.Add("limit", strconv.Itoa(qlimit))
		q.Add("offset", strconv.Itoa(qoffset))
		req.URL.RawQuery = q.Encode()

		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body, err := json.Marshal(keywords)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), w.Body.String())

		keywordRepo.AssertExpectations(t)
	})

}
