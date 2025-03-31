package nodeitem_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"pan/lib/web"
	"path/filepath"
	"testing"

	nodeitem "pan/features/extfs/node_item"
	mockedNodeItem "pan/mocks/pan/features/extfs/node_item"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeFileStreamController(t *testing.T) {

	setup := func() (web.WebApp, *nodeitem.NodeFileStreamController) {
		app := web.NewWebApp()
		ctrl := &nodeitem.NodeFileStreamController{}
		ctrl.SetupToWeb(app)

		ctrl.NodeFilePathService = &nodeitem.NodeFilePathService{}
		return app, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /node-items/:id/_stream/*filepath", func(t *testing.T) {
		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-item-streams")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.json"
		fullpath := filepath.Join(filePath, filename)
		fileBytes := []byte("{\"name\":\"this is test json\"}")
		err = os.WriteFile(fullpath, fileBytes, 0644)
		assert.Nil(t, err)

		nodeItemService := &mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.NodeFilePathService.NodeItemService = nodeItemService

		itemId := uint(1)
		var nodeItem nodeitem.NodeItem
		nodeItem.ID = itemId
		nodeItem.FilePath = filePath
		nodeItem.FileType = nodeitem.FileTypeFolder
		nodeItem.Available = true

		nodeItemService.On("Select", itemId).Once().Return(nodeItem, nil)
		nodeItemService.On("IsNotExist", mock.Anything).Once().Return(false)

		url := fmt.Sprintf("/node-items/%d/_stream/%s", itemId, filename)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		resBodyBytes := w.Body.Bytes()
		assert.Equal(t, fileBytes, resBodyBytes)

	})

	t.Run("GET /node-items/:id/_stream/*filepath with Range", func(t *testing.T) {
		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-item-streams")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.json"
		fullpath := filepath.Join(filePath, filename)
		fileBytes := []byte("{\"name\":\"this is test json\"}")
		err = os.WriteFile(fullpath, fileBytes, 0644)
		assert.Nil(t, err)

		nodeItemService := &mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.NodeFilePathService.NodeItemService = nodeItemService

		itemId := uint(1)
		var nodeItem nodeitem.NodeItem
		nodeItem.ID = itemId
		nodeItem.FilePath = filePath
		nodeItem.FileType = nodeitem.FileTypeFolder
		nodeItem.Available = true

		nodeItemService.On("Select", itemId).Once().Return(nodeItem, nil)
		nodeItemService.On("IsNotExist", mock.Anything).Once().Return(false)

		url := fmt.Sprintf("/node-items/%d/_stream/%s", itemId, filename)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Range", "bytes=0-2000")
		web.ServeHTTP(w, req)

		assert.Equal(t, 206, w.Code)

		resBodyBytes := w.Body.Bytes()
		assert.Equal(t, fileBytes, resBodyBytes)

	})

	t.Run("GET /node-items/:id/_stream/*filepath Failed with Folder", func(t *testing.T) {
		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-item-streams")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.json"
		fullpath := filepath.Join(filePath, filename)
		fileBytes := []byte("{\"name\":\"this is test json\"}")
		err = os.WriteFile(fullpath, fileBytes, 0644)
		assert.Nil(t, err)

		nodeItemService := &mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.NodeFilePathService.NodeItemService = nodeItemService

		itemId := uint(1)
		var nodeItem nodeitem.NodeItem
		nodeItem.ID = itemId
		nodeItem.FilePath = filePath
		nodeItem.FileType = nodeitem.FileTypeFolder
		nodeItem.Available = true

		nodeItemService.On("Select", itemId).Once().Return(nodeItem, nil)
		nodeItemService.On("IsNotExist", mock.Anything).Once().Return(false)

		url := fmt.Sprintf("/node-items/%d/_stream/", itemId)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		web.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)

	})

}
