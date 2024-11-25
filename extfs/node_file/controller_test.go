package nodefile_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"pan/app/web"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	nodefile "pan/extfs/node_file"
	nodeitem "pan/extfs/node_item"
	mockedNodeItem "pan/mocks/pan/extfs/node_item"
)

func TestNodeFileController(t *testing.T) {

	setup := func() (web.WebApp, *nodefile.NodeFileController) {
		app := web.NewWebApp()
		ctrl := &nodefile.NodeFileController{}
		ctrl.SetupToWeb(app)

		ctrl.NodeFileService = &nodefile.NodeFileService{}
		return app, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /node-files?itemId=?", func(t *testing.T) {

		web, ctrl := setup()

		filePath, err := setupTemp("extfs-node-files")
		assert.Nil(t, err)
		defer teardownTemp(filePath)

		filename := "test.txt"
		fullpath := path.Join(filePath, filename)
		err = os.WriteFile(fullpath, []byte("hello"), 0644)
		assert.Nil(t, err)
		fileStat, err := os.Stat(fullpath)
		assert.Nil(t, err)

		nodeItemService := &mockedNodeItem.MockNodeItemInternalService{}
		defer nodeItemService.AssertExpectations(t)
		ctrl.NodeFileService.NodeItemService = nodeItemService

		itemId := uint(1)
		var nodeItem nodeitem.NodeItem
		nodeItem.ID = itemId
		nodeItem.FilePath = filePath
		nodeItem.FileType = nodeitem.FileTypeFolder
		nodeItem.Available = true

		nodeItemService.On("Select", itemId).Once().Return(nodeItem, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/node-files", nil)
		q := req.URL.Query()
		q.Add("itemId", strconv.Itoa(int(itemId)))
		req.URL.RawQuery = q.Encode()
		web.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var results []nodefile.NodeFile
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
