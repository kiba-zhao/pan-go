package services

import (
	"encoding/base64"
	"os"
	"pan/app/constant"
	"pan/app/models"
	"path"
	"sync"
)

type DiskFileService struct {
	rw   sync.RWMutex
	Root *models.DiskFile
}

func (s *DiskFileService) Search(conditions models.DiskFileSearchCondition) (total int64, items []models.DiskFile, err error) {

	if conditions.FilePath != "" || conditions.ParentPath == "" {

		var item models.DiskFile
		var itemErr error
		if conditions.FilePath != "" {
			item, itemErr = s.SelectWithFilePath(conditions.FilePath)
		} else {
			item, itemErr = s.SelectRoot()
		}

		if itemErr != nil {
			err = itemErr
			return
		}

		if conditions.FileType != "" && conditions.FileType != item.FileType {
			return
		}

		if conditions.ParentPath != "" && conditions.ParentPath != item.ParentPath {
			err = constant.ErrConflict
			return
		}

		total = 1
		items = append(items, item)
		return
	}

	dirs, err := os.ReadDir(conditions.ParentPath)
	if err != nil {
		return
	}

	total = int64(len(dirs))
	if total <= 0 {
		return
	}
	for _, dir := range dirs {

		info, err := dir.Info()
		if err != nil {
			break
		}
		if conditions.FileType != "" && conditions.FileType != getFileType(info.IsDir()) {
			continue
		}
		filePath := path.Join(conditions.ParentPath, dir.Name())
		items = append(items, models.DiskFile{
			ID:         encodeFilePath(filePath),
			Name:       dir.Name(),
			FilePath:   filePath,
			ParentPath: conditions.ParentPath,
			FileType:   getFileType(info.IsDir()),
			UpdatedAt:  info.ModTime(),
		})
	}

	return
}

func (s *DiskFileService) SelectRoot() (item models.DiskFile, err error) {

	s.rw.RLock()
	if s.Root != nil {
		defer s.rw.RUnlock()
		item = *s.Root
		return
	}
	s.rw.RUnlock()

	s.rw.Lock()
	defer s.rw.Unlock()
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	item, err = s.SelectWithFilePath(rootPath)
	if err == nil {
		s.Root = &item
	}
	return
}

func (s *DiskFileService) SelectWithFilePath(filePath string) (item models.DiskFile, err error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return
	}

	item = models.DiskFile{
		ID:         encodeFilePath(filePath),
		Name:       stat.Name(),
		FilePath:   filePath,
		ParentPath: path.Dir(filePath),
		FileType:   getFileType(stat.IsDir()),
		UpdatedAt:  stat.ModTime(),
	}
	return
}

func getFileType(isDir bool) string {
	if isDir {
		return models.FILETYPE_FOLDER
	}
	return models.FILETYPE_FILE
}

func encodeFilePath(filePath string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(filePath))
	return encoded
}
