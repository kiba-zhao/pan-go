package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"pan/app/net"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"strconv"
	"testing"

	mockedRepo "pan/mocks/pan/extfs/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTargetFile(t *testing.T) {

	setup := func() (web net.WebApp, ctrl *controllers.TargetFileController) {
		ctrl = new(controllers.TargetFileController)
		web = net.NewWebApp()
		ctrl.SetupToWeb(web)

		ctrl.TargetFileService = &services.TargetFileService{}
		return web, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /target-files", func(t *testing.T) {
		web, ctrl := setup()

		targetFileRepo := new(mockedRepo.MockTargetFileRepository)
		defer targetFileRepo.AssertExpectations(t)
		ctrl.TargetFileService.TargetFileRepo = targetFileRepo

		enabled := false
		var target models.Target
		target.ID = 1
		target.Enabled = &enabled
		total := int64(10)
		targetFiles := []models.TargetFile{
			{ID: 1, FilePath: "/path_a", Target: target},
			{ID: 2, FilePath: "/path_b", Target: target},
		}
		targetFileRepo.On("Search", models.TargetFileSearchCondition{}, true).Once().Return(total, targetFiles, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/target-files", nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.TargetFile
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Len(t, results, len(targetFiles))
		for idx, item := range results {
			targetFile := targetFiles[idx]
			assert.Equal(t, targetFile.ID, item.ID)
			assert.Equal(t, targetFile.FilePath, item.FilePath)
			assert.False(t, item.Available)
		}
	})

	t.Run("GET /target-files?target_id=&keyword=&available=", func(t *testing.T) {
		web, ctrl := setup()

		targetFileRepo := new(mockedRepo.MockTargetFileRepository)
		defer targetFileRepo.AssertExpectations(t)
		ctrl.TargetFileService.TargetFileRepo = targetFileRepo

		var condition models.TargetFileSearchCondition
		condition.TargetID = 1
		condition.Keyword = "keyword"
		condition.RangeStart = 0
		condition.RangeEnd = 12
		condition.SortField = "name"
		condition.SortOrder = "desc"

		enabled := false
		var target models.Target
		target.ID = 1
		target.Enabled = &enabled
		total := int64(10)
		targetFiles := []models.TargetFile{
			{ID: 1, FilePath: "/path_a", Target: target},
			{ID: 2, FilePath: "/path_b", Target: target},
		}
		targetFileRepo.On("Search", mock.AnythingOfType("models.TargetFileSearchCondition"), true).Once().Return(total, targetFiles, nil).Run(func(args mock.Arguments) {
			condition_ := args.Get(0).(models.TargetFileSearchCondition)

			assert.Equal(t, condition.TargetID, condition_.TargetID)
			assert.Equal(t, condition.Keyword, condition_.Keyword)

			assert.Equal(t, condition.RangeStart, condition_.RangeStart)
			assert.Equal(t, condition.RangeEnd, condition_.RangeEnd)

			assert.Equal(t, condition.SortField, condition_.SortField)
			assert.Equal(t, condition.SortOrder, condition_.SortOrder)

		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/target-files", nil)
		q := req.URL.Query()
		q.Add("targetId", strconv.FormatUint(uint64(condition.TargetID), 10))
		q.Add("q", condition.Keyword)
		q.Add("_start", strconv.Itoa(condition.RangeStart))
		q.Add("_end", strconv.Itoa(condition.RangeEnd))
		q.Add("_sort", condition.SortField)
		q.Add("_order", condition.SortOrder)
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(net.CountHeaderName))
		var results []models.TargetFile
		err := json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Len(t, results, len(targetFiles))
		for idx, item := range results {
			targetFile := targetFiles[idx]
			assert.Equal(t, targetFile.ID, item.ID)
			assert.Equal(t, targetFile.FilePath, item.FilePath)
			assert.False(t, item.Available)
		}
	})

	t.Run("GET /target-files/:id", func(t *testing.T) {

		web, ctrl := setup()

		targetFileRepo := new(mockedRepo.MockTargetFileRepository)
		defer targetFileRepo.AssertExpectations(t)
		ctrl.TargetFileService.TargetFileRepo = targetFileRepo

		var targetFile models.TargetFile
		targetFile.ID = 1
		targetFile.Available = true
		targetFile.FilePath = os.TempDir()

		enabled := true
		var target models.Target
		target.ID = 123
		target.HashCode = "hash"
		target.FilePath = os.TempDir()
		target.Enabled = &enabled

		targetFile.TargetID = target.ID
		targetFile.Target = target
		targetFileRepo.On("Select", targetFile.ID, true).Once().Return(targetFile, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/target-files/%d", targetFile.ID)
		req, _ := http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var result models.TargetFile
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, targetFile.Available, result.Available)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.TargetID, result.TargetID)

		enabled = false
		targetFile.Target = target
		targetFileRepo.On("Select", targetFile.ID, true).Once().Return(targetFile, nil)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, false, result.Available)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.TargetID, result.TargetID)

		invaildPath, err := setupTemp("extfs-target-files-test")
		assert.Nil(t, err)
		err = teardownTemp(invaildPath)
		assert.Nil(t, err)

		enabled = true
		target.FilePath = invaildPath
		targetFile.Target = target
		targetFileRepo.On("Select", targetFile.ID, true).Once().Return(targetFile, nil)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, false, result.Available)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.TargetID, result.TargetID)

		enabled = true
		target.FilePath = os.TempDir()
		targetFile.Target = target
		targetFile.FilePath = invaildPath
		targetFileRepo.On("Select", targetFile.ID, true).Once().Return(targetFile, nil)

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &result)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, false, result.Available)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.TargetID, result.TargetID)
	})
}
