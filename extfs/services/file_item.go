package services

import (
	"encoding/base64"
	"os"
	appConstant "pan/app/constant"
	"pan/extfs/constant"
	"pan/extfs/models"
	"path"
	"strconv"
	"strings"
)

type FileItemInternalService interface {
	TraverseWithCondition(traverseFn func(item models.FileItem) error, conditions models.FileItemSearchCondition) error
}

type FileItemService struct {
	NodeItemService NodeItemInternalService
}

func (s *FileItemService) Search(conditions models.FileItemSearchCondition) (int64, []models.FileItem, error) {

	items := make([]models.FileItem, 0)
	err := s.TraverseWithCondition(func(item models.FileItem) error {
		items = append(items, item)
		return nil
	}, conditions)
	return int64(len(items)), items, err
}

func (s *FileItemService) TraverseWithCondition(traverseFn func(item models.FileItem) error, conditions models.FileItemSearchCondition) error {
	nodeItem, err := s.NodeItemService.Select(conditions.ItemID)
	if err != nil {
		return err
	}

	if !nodeItem.Available || nodeItem.FileType != FileTypeFolder {
		return appConstant.ErrUnavailable
	}

	filePath := nodeItem.FilePath
	if conditions.ParentPath != nil {
		filePath = path.Join(filePath, *conditions.ParentPath)
	}

	files, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		var item models.FileItem

		item.ItemID = nodeItem.ID
		item.Name = file.Name()
		item.ID = generateFileItemID(item.ItemID, item.Name)

		if file.IsDir() {
			item.FileType = constant.FileTypeDir
		} else {
			item.FileType = constant.FileTypeFile
		}

		if conditions.ParentPath == nil {
			item.FilePath = item.Name
		} else {
			item.ParentPath = *conditions.ParentPath
			item.FilePath = path.Join(*conditions.ParentPath, item.Name)

		}

		item.Available = true
		info, infoErr := file.Info()
		if infoErr == nil {
			item.UpdatedAt = info.ModTime()
			item.Size = info.Size()
		} else {
			item.Available = false
		}
		err = traverseFn(item)
		if err != nil {
			return err
		}
	}
	return err
}

func generateFileItemID(itemId uint, filePath string) string {
	idStr := strings.Join([]string{strconv.FormatUint(uint64(itemId), 10), filePath}, constant.FileItemSep)
	return base64.StdEncoding.EncodeToString([]byte(idStr))
}
