package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"pan/app/net"
	"pan/extfs/controllers"
	"pan/extfs/models"
	"pan/extfs/services"
	"path"
	"strconv"
	"testing"

	mockedServices "pan/mocks/pan/extfs/services"

	"github.com/stretchr/testify/assert"
)

func TestFileItemController(t *testing.T) {

	setup := func() (net.WebApp, *controllers.FileItemController) {
		app := net.NewWebApp()
		ctrl := &controllers.FileItemController{}
		ctrl.SetupToWeb(app)

		ctrl.FileItemService = &services.FileItemService{}
		return app, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /file-items?itemId=?", func(t *testing.T) {

		web, ctrl := setup()

		filepath, err := setupTemp("extfs-file-items")
		assert.Nil(t, err)
		defer teardownTemp(filepath)

		filename := "test.txt"
		fullpath := path.Join(filepath, filename)
		err = os.WriteFile(fullpath, []byte("hello"), 0644)
		assert.Nil(t, err)
		fileStat, err := os.Stat(fullpath)
		assert.Nil(t, err)

		nodeItemService := &mockedServices.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.FileItemService.NodeItemService = nodeItemService

		itemId := uint(1)
		var nodeItem models.NodeItem
		nodeItem.ID = itemId
		nodeItem.FilePath = filepath
		nodeItem.FileType = services.FileTypeFolder
		nodeItem.Available = true

		nodeItemService.On("Select", itemId).Once().Return(nodeItem, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/file-items", nil)
		q := req.URL.Query()
		q.Add("itemId", strconv.Itoa(int(itemId)))
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []models.FileItem
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, filename, results[0].Name)
		assert.Equal(t, fileStat.Size(), results[0].Size)
		assert.Equal(t, filename, results[0].FilePath)
		assert.Equal(t, services.FileTypeFile, results[0].FileType)
		assert.True(t, results[0].Available)
	})
}
