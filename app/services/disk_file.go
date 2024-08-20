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

	if conditions.FilePath != "" || conditions.Parent == "" {

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

		if conditions.Parent != "" && conditions.Parent != item.Parent {
			err = constant.ErrConflict
			return
		}

		total = 1
		items = append(items, item)
		return
	}

	dirs, err := os.ReadDir(conditions.Parent)
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
		filepath := path.Join(conditions.Parent, dir.Name())
		items = append(items, models.DiskFile{
			ID:        encodeFilePath(filepath),
			Name:      dir.Name(),
			FilePath:  filepath,
			Parent:    conditions.Parent,
			FileType:  getFileType(info.IsDir()),
			UpdatedAt: info.ModTime(),
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

func (s *DiskFileService) SelectWithFilePath(filepath string) (item models.DiskFile, err error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return
	}

	item = models.DiskFile{
		ID:        encodeFilePath(filepath),
		Name:      stat.Name(),
		FilePath:  filepath,
		Parent:    path.Dir(filepath),
		FileType:  getFileType(stat.IsDir()),
		UpdatedAt: stat.ModTime(),
	}
	return
}

func getFileType(isDir bool) string {
	if isDir {
		return models.FILETYPE_FOLDER
	}
	return models.FILETYPE_FILE
}

func encodeFilePath(filepath string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(filepath))
	return encoded
}
