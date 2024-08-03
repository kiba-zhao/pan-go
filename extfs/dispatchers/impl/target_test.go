package impl_test

import (
	"os"
	"pan/extfs/dispatchers/impl"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
	"path"
	"sync"
	"testing"

	mockedRepo "pan/mocks/pan/extfs/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestTarget(t *testing.T) {

	setup := func() (dispatcher *impl.TargetDispatcher) {
		dispatcher = impl.NewTargetDispatcher()

		dispatcher.TargetService = &services.TargetService{}
		dispatcher.TargetService.TargetFileService = &services.TargetFileService{}
		return
	}

	setupTemp := func(name string) (string, error) {
		dir, err := os.MkdirTemp(os.TempDir(), name)
		return dir, err
	}

	teardownTemp := func(dir string) error {
		return os.RemoveAll(dir)
	}

	t.Run("Scan", func(t *testing.T) {
		dispatcher := setup()

		root, err := setupTemp("extfs-target-test")
		assert.Nil(t, err)
		defer teardownTemp(root)

		folderName := "folder1"
		folderPath := path.Join(root, folderName)
		err = os.MkdirAll(folderPath, 0755)
		assert.Nil(t, err)
		fileName := "file1.txt"
		filePath := path.Join(root, fileName)
		os.WriteFile(filePath, []byte("hello"), 0644)

		enabled := true
		target := models.Target{}
		target.ID = uint(123)
		target.FilePath = ""
		target.Name = "Target A"
		target.Enabled = &enabled
		target.Available = true

		targetRepo := new(mockedRepo.MockTargetRepository)
		dispatcher.TargetService.TargetRepo = targetRepo
		defer targetRepo.AssertExpectations(t)

		newTarget := models.Target{}
		newTarget.ID = target.ID
		newTarget.HashCode = "hash code"
		newTarget.FilePath = root
		newTarget.Name = "Target B"
		newTarget.Enabled = &enabled
		newTarget.Available = true
		targetRepo.On("Select", target.ID, mock.Anything).Return(newTarget, nil)

		var wg sync.WaitGroup
		wg.Add(1)

		targetFileRepo := new(mockedRepo.MockTargetFileRepository)
		dispatcher.TargetService.TargetFileService.TargetFileRepo = targetFileRepo
		defer targetFileRepo.AssertExpectations(t)
		targetFileRepo.On("TraverseByTargetId", mock.Anything, newTarget.ID).Once().Return(nil)
		targetFileRepo.On("SelectByFilePathAndTargetId", filePath, newTarget.ID, mock.AnythingOfType("string"), false).Once().Return(models.TargetFile{}, errors.ErrNotFound)

		var saveTargetFile models.TargetFile
		targetFileRepo.On("Save", mock.Anything).Once().Return(saveTargetFile, nil).Run(func(args mock.Arguments) {
			defer wg.Done()
			targetFile := args.Get(0).(models.TargetFile)

			saveTargetFile.TargetID = targetFile.TargetID
			saveTargetFile.FilePath = targetFile.FilePath
			saveTargetFile.TargetHashCode = targetFile.TargetHashCode
			saveTargetFile.HashCode = targetFile.TargetHashCode

			assert.Equal(t, targetFile.TargetID, newTarget.ID)
			assert.Equal(t, targetFile.FilePath, filePath)
			assert.Equal(t, targetFile.TargetHashCode, newTarget.HashCode)

		})

		err = dispatcher.Scan(target)
		assert.Nil(t, err)
		if err != nil {
			wg.Done()
		}

		wg.Wait()
	})

	t.Run("Clean", func(t *testing.T) {
		dispatcher := setup()

		target := models.Target{}
		target.ID = uint(123)
		target.DeletedAt = gorm.DeletedAt{Valid: true}

		var wg sync.WaitGroup
		wg.Add(1)

		targetRepo := new(mockedRepo.MockTargetRepository)
		dispatcher.TargetService.TargetRepo = targetRepo
		defer targetRepo.AssertExpectations(t)
		targetRepo.On("Select", target.ID, mock.Anything).Return(target, nil)

		targetFileRepo := new(mockedRepo.MockTargetFileRepository)
		dispatcher.TargetService.TargetFileService.TargetFileRepo = targetFileRepo
		defer targetFileRepo.AssertExpectations(t)
		targetFileRepo.On("DeleteByTargetId", target.ID).Once().Return(nil).Run(func(args mock.Arguments) {
			defer wg.Done()
		})

		err := dispatcher.Clean(target)
		assert.Nil(t, err)
		wg.Wait()
	})
}
