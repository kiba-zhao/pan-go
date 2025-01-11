package nodesearchfile_test

import (
	"os"
	nodesearchfile "pan/extfs/node_search_file"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFSSearcherFS(t *testing.T) {

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("WalkRoot", func(t *testing.T) {

		dir, err := setupTemp("extfs-nodesearchfile-fs-test")
		assert.Nil(t, err)
		defer teardownTemp(dir)

		folderName := "folder1"
		folderPath := path.Join(dir, folderName)
		err = os.MkdirAll(folderPath, 0755)
		assert.Nil(t, err)

		filderFileName := "subfile.txt"
		filderFilePath := path.Join(folderPath, filderFileName)
		os.WriteFile(filderFilePath, []byte("hello world"), 0644)

		fileName := "file1.txt"
		filePath := path.Join(dir, fileName)
		os.WriteFile(filePath, []byte("hello"), 0644)

		filePaths, err := nodesearchfile.WalkRoot(dir)
		assert.Nil(t, err)

		count := 0
		for filePath_ := range filePaths {
			assert.Contains(t, []string{dir, folderPath, filePath, filderFilePath}, filePath_)
			count++
		}
		assert.Equal(t, 4, count)
	})

}
