package diskfile_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"pan/app/web"
	"path/filepath"
	"strconv"
	"testing"

	diskfile "pan/app/disk_file"

	"github.com/stretchr/testify/assert"
)

func TestDiskFileController(t *testing.T) {

	setup := func() (web.WebApp, *diskfile.DiskFileController) {
		ctrl := new(diskfile.DiskFileController)
		webApp := web.NewWebApp()
		ctrl.SetupToWeb(webApp)

		ctrl.DiskFileService = &diskfile.DiskFileService{}
		return webApp, ctrl
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("GET /disk-files?parent=", func(t *testing.T) {
		webApp, ctrl := setup()

		var root diskfile.DiskFile
		root.Name = "root"
		root.FilePath = "/"
		root.ParentPath = "/parent"
		root.FileType = diskfile.FILETYPE_FOLDER
		ctrl.DiskFileService.Root = &root

		total := int64(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/disk-files?parent=", nil)
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))

		var rows []diskfile.DiskFile
		err := json.Unmarshal(w.Body.Bytes(), &rows)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(rows))

		assert.Equal(t, root.Name, rows[0].Name)
		assert.Equal(t, root.FilePath, rows[0].FilePath)
		assert.Equal(t, root.ParentPath, rows[0].ParentPath)
		assert.Equal(t, root.FileType, rows[0].FileType)
	})

	t.Run("GET /disk-files?parentPath=/path", func(t *testing.T) {
		webApp, _ := setup()

		parent, err := setupTemp("extfs-disk-files-test")
		assert.Nil(t, err)
		defer teardownTemp(parent)

		folderName := "folder1"
		folderPath := filepath.Join(parent, folderName)
		err = os.MkdirAll(folderPath, 0755)
		assert.Nil(t, err)
		fileName := "file1.txt"
		filePath := filepath.Join(parent, fileName)
		os.WriteFile(filePath, []byte("hello"), 0644)

		total := int64(2)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/disk-files", nil)
		q := req.URL.Query()
		q.Add("parentPath", parent)
		req.URL.RawQuery = q.Encode()
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))
		var rows []diskfile.DiskFile
		err = json.Unmarshal(w.Body.Bytes(), &rows)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(rows))

		fileItem := rows[0]
		assert.Equal(t, fileName, fileItem.Name)
		assert.Equal(t, filePath, fileItem.FilePath)
		assert.Equal(t, parent, fileItem.ParentPath)
		assert.Equal(t, diskfile.FILETYPE_FILE, fileItem.FileType)

		dirItem := rows[1]
		assert.Equal(t, folderName, dirItem.Name)
		assert.Equal(t, folderPath, dirItem.FilePath)
		assert.Equal(t, parent, dirItem.ParentPath)
		assert.Equal(t, diskfile.FILETYPE_FOLDER, dirItem.FileType)

	})

	t.Run("GET /disk-files?filePath=/path", func(t *testing.T) {
		webApp, _ := setup()

		parent, err := setupTemp("extfs-disk-files-test")
		assert.Nil(t, err)
		defer teardownTemp(parent)

		folderName := "folder2"
		folderPath := filepath.Join(parent, folderName)
		err = os.MkdirAll(folderPath, 0755)
		assert.Nil(t, err)

		total := int64(1)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/disk-files", nil)
		q := req.URL.Query()
		q.Add("filePath", folderPath)
		req.URL.RawQuery = q.Encode()
		webApp.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, strconv.FormatInt(total, 10), w.Header().Get(web.CountHeaderName))
		var rows []diskfile.DiskFile
		err = json.Unmarshal(w.Body.Bytes(), &rows)
		assert.Nil(t, err)

		assert.Equal(t, 1, len(rows))

		dirItem := rows[0]
		assert.Equal(t, folderName, dirItem.Name)
		assert.Equal(t, folderPath, dirItem.FilePath)
		assert.Equal(t, parent, dirItem.ParentPath)
		assert.Equal(t, diskfile.FILETYPE_FOLDER, dirItem.FileType)

	})
}
