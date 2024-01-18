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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestTags ...
func TestTags(t *testing.T) {
	setup := func(app *web.App, repo *repositories.MockTagRepository, limit int) {
		ctrl := new(controllers.TagController)
		ctrl.TagService = new(services.TagService)
		ctrl.TagService.TagRepo = repo
		ctrl.DefaultLimit = limit

		ctrl.Init(app)
	}

	t.Run("GET /extfs/tags", func(t *testing.T) {
		limit := 101
		tagRepo := new(mocked.MockTagRepository)
		app := web.NewApp()
		setup(app, tagRepo, 101)

		tag := new(models.Tag)
		tag.Name = "tag name"
		tag.Owner = uuid.New().String()
		tags := []models.Tag{*tag}

		condition := new(models.TagFindCondition)
		condition.Limit = limit
		tagRepo.On("Find", condition).Once().Return(tags, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/extfs/tags", nil)
		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body, err := json.Marshal(tags)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), w.Body.String())

		tagRepo.AssertExpectations(t)
	})

	t.Run("GET /extfs/tags with special name", func(t *testing.T) {

		limit := 100
		tagRepo := new(mocked.MockTagRepository)
		app := web.NewApp()
		setup(app, tagRepo, limit)

		qname := "ta"
		qlimit := 10
		qoffset := 0
		tag := new(models.Tag)
		tag.Name = "tag name"
		tag.Owner = uuid.New().String()
		tags := []models.Tag{*tag}

		condition := new(models.TagFindCondition)
		condition.Name = qname
		condition.Limit = qlimit
		condition.Offset = qoffset
		tagRepo.On("Find", condition).Once().Return(tags, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/extfs/tags", nil)
		q := req.URL.Query()
		q.Add("name", qname)
		q.Add("limit", strconv.Itoa(qlimit))
		q.Add("offset", strconv.Itoa(qoffset))
		req.URL.RawQuery = q.Encode()

		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body, err := json.Marshal(tags)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), w.Body.String())

		tagRepo.AssertExpectations(t)
	})

	t.Run("GET /extfs/tags with invalid owner", func(t *testing.T) {

		limit := 100
		tagRepo := new(mocked.MockTagRepository)
		app := web.NewApp()
		setup(app, tagRepo, limit)

		qname := "ta"
		qlimit := 10
		qoffset := 0
		tag := new(models.Tag)
		tag.Name = "tag name"
		tag.Owner = uuid.New().String()
		tags := []models.Tag{*tag}

		tagRepo.On("Find", (*models.TagFindCondition)(nil)).Once().Return(tags, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/extfs/tags", nil)
		q := req.URL.Query()
		q.Add("owner", "123")
		q.Add("name", qname)
		q.Add("limit", strconv.Itoa(qlimit))
		q.Add("offset", strconv.Itoa(qoffset))
		req.URL.RawQuery = q.Encode()

		app.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		body, err := json.Marshal(tags)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, string(body), w.Body.String())

		tagRepo.AssertExpectations(t)
	})

}
