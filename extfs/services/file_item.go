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

type FileItemService struct {
	NodeItemService NodeItemInternalService
}

func (s *FileItemService) Search(conditions models.FileItemSearchCondition) (int64, []models.FileItem, error) {
	nodeItem, err := s.NodeItemService.Select(conditions.ItemID)
	if err != nil {
		return 0, nil, err
	}

	if !nodeItem.Available || nodeItem.FileType != FileTypeFolder {
		return 0, nil, appConstant.ErrUnavailable
	}

	filePath := nodeItem.FilePath
	if conditions.ParentPath != nil {
		filePath = path.Join(filePath, *conditions.ParentPath)
	}

	files, err := os.ReadDir(filePath)
	if err != nil {
		return 0, nil, err
	}

	items := make([]models.FileItem, 0)
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
		items = append(items, item)
	}
	return int64(len(items)), items, nil

}

func generateFileItemID(itemId uint, filePath string) string {
	idStr := strings.Join([]string{strconv.FormatUint(uint64(itemId), 10), filePath}, constant.FileItemSep)
	return base64.StdEncoding.EncodeToString([]byte(idStr))
}
