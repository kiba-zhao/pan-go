package diskfile

import (
	"encoding/base64"
	"errors"
	"os"
	"path"
	"sync"
)

var ErrDiskFileParentPathConflict = errors.New("diskfile.DiskFileService Error: Parent Path Conflict")

type DiskFileService struct {
	rw   sync.RWMutex
	Root *DiskFile
}

func (s *DiskFileService) Search(conditions DiskFileSearchCondition) (total int64, items []DiskFile, err error) {

	if conditions.FilePath != "" || conditions.ParentPath == "" {

		var item DiskFile
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
			err = ErrDiskFileParentPathConflict
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
		items = append(items, DiskFile{
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

func (s *DiskFileService) SelectRoot() (item DiskFile, err error) {

	s.rw.RLock()
	if s.Root != nil {
		defer s.rw.RUnlock()
		item = *s.Root
		return
	}
	s.rw.RUnlock()

	s.rw.Lock()
	defer s.rw.Unlock()
	rootPath, err := os.UserHomeDir()
	if err != nil {
		return
	}
	item, err = s.SelectWithFilePath(rootPath)
	if err == nil {
		s.Root = &item
	}
	return
}

func (s *DiskFileService) SelectWithFilePath(filePath string) (item DiskFile, err error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return
	}

	item = DiskFile{
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
		return FILETYPE_FOLDER
	}
	return FILETYPE_FILE
}

func encodeFilePath(filePath string) string {
	encoded := base64.RawURLEncoding.EncodeToString([]byte(filePath))
	return encoded
}
