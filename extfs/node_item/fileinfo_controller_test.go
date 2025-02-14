package nodeitem_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"pan/app/web"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	nodeitem "pan/extfs/node_item"
	mockedNodeItem "pan/mocks/pan/extfs/node_item"
)

func TestNodeFileInfoController(t *testing.T) {

	setup := func() (web.WebApp, *nodeitem.NodeFileInfoController) {
		app := web.NewWebApp()
		ctrl := &nodeitem.NodeFileInfoController{}
		ctrl.SetupToWeb(app)

		ctrl.NodeFileInfoService = &nodeitem.NodeFileInfoService{}
		return app, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /node-items/:id/file/:filepath", func(t *testing.T) {
		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-item-files")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.txt"
		fullpath := filepath.Join(filePath, filename)
		err = os.WriteFile(fullpath, []byte("hello"), 0644)
		assert.Nil(t, err)
		fileStat, err := os.Stat(fullpath)
		assert.Nil(t, err)

		nodeFilePathService := &mockedNodeItem.MockNodeFilePathInternalService{}
		defer nodeFilePathService.AssertExpectations(t)
		ctrl.NodeFileInfoService.NodeFilePathService = nodeFilePathService

		itemId := uint(1)
		nodeFilePathService.On("Select", itemId, filename).Once().Return(fullpath, nil)
		nodeFilePathService.On("IsNotExist", mock.Anything).Once().Return(false)

		url := fmt.Sprintf("/node-items/%d/file/%s", itemId, filename)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var nodeFile nodeitem.NodeFileInfo
		err = json.Unmarshal(w.Body.Bytes(), &nodeFile)
		assert.Nil(t, err)
		assert.Equal(t, itemId, nodeFile.ItemID)
		assert.Equal(t, filename, nodeFile.Name)
		assert.Equal(t, nodeitem.FileTypeFile, nodeFile.FileType)
		assert.Equal(t, filename, nodeFile.FilePath)
		assert.Equal(t, fileStat.Size(), nodeFile.Size)

	})

	t.Run("GET /node-items/:id/files/:filepath", func(t *testing.T) {

		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-item-files")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.txt"
		fullpath := filepath.Join(filePath, filename)
		err = os.WriteFile(fullpath, []byte("hello"), 0644)
		assert.Nil(t, err)
		fileStat, err := os.Stat(fullpath)
		assert.Nil(t, err)

		nodeFilePathService := &mockedNodeItem.MockNodeFilePathInternalService{}
		defer nodeFilePathService.AssertExpectations(t)
		ctrl.NodeFileInfoService.NodeFilePathService = nodeFilePathService

		itemId := uint(1)

		nodeFilePathService.On("IsNotExist", mock.Anything).Once().Return(false)
		nodeFilePathService.On("Select", itemId, "").Once().Return(filePath, nil)

		w := httptest.NewRecorder()
		url := fmt.Sprintf("/node-items/%d/files/", itemId)
		req := httptest.NewRequest("GET", url, nil)
		q := req.URL.Query()
		q.Add("itemId", strconv.Itoa(int(itemId)))
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []nodeitem.NodeFileInfo
		err = json.Unmarshal(w.Body.Bytes(), &results)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(results))

		assert.Equal(t, filename, results[0].Name)
		assert.Equal(t, fileStat.Size(), results[0].Size)
		assert.Equal(t, filename, results[0].FilePath)
		assert.Equal(t, nodeitem.FileTypeFile, results[0].FileType)
		assert.True(t, results[0].Available)
	})
}
