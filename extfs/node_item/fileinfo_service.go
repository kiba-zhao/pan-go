package nodeitem

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

var ErrNodeFileInfoUnavailable = errors.New("nodeitem.NodeFileInfoService Error: Unavailable")
var ErrNodeFileInfoInvalidNodeItem = errors.New("nodeitem.NodeFileInfoService Error: Invalid Node Item")
var ErrNodeFileInfoInvalidID = errors.New("nodeitem.NodeFileInfoService Error: Invalid Node File ID")

type NodeFileInfoInternalService interface {
	TraverseWithCondition(func(fileInfo NodeFileInfo) error, NodeFileInfoSearchCondition) error
	Select(id uint, filePath string) (NodeFileInfo, error)
	IsNotExist(err error) bool
}

type NodeFileInfoService struct {
	NodeFilePathService NodeFilePathInternalService
}

func (s *NodeFileInfoService) IsNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if errors.Is(err, ErrNodeFileInfoUnavailable) {
		return true
	}
	return s.NodeFilePathService.IsNotExist(err)
}

func (s *NodeFileInfoService) Select(id uint, filePath string) (NodeFileInfo, error) {
	realFilePath, err := s.NodeFilePathService.Select(id, filePath)
	if err != nil {
		return NodeFileInfo{}, err
	}

	fileStat, err := os.Stat(realFilePath)
	if err != nil {
		return NodeFileInfo{}, err
	}

	fileInfo := generateFileInfo(id, filePath, fileStat)
	return fileInfo, err
}

func (s *NodeFileInfoService) Search(id uint, parentPath string) (int64, []NodeFileInfo, error) {

	var conditions NodeFileInfoSearchCondition
	conditions.ItemID = id
	conditions.ParentPath = parentPath

	items := make([]NodeFileInfo, 0)
	err := s.TraverseWithCondition(func(item NodeFileInfo) error {
		items = append(items, item)
		return nil
	}, conditions)
	return int64(len(items)), items, err
}

func (s *NodeFileInfoService) TraverseWithCondition(traverseFn func(item NodeFileInfo) error, conditions NodeFileInfoSearchCondition) error {

	realFilePath, err := s.NodeFilePathService.Select(conditions.ItemID, conditions.ParentPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(realFilePath)
	if err != nil {
		return err
	}

	for _, file := range files {

		filePath := file.Name()
		if len(conditions.ParentPath) > 0 {
			filePath = path.Join(conditions.ParentPath, filePath)
		}

		var fileInfo NodeFileInfo
		info, infoErr := file.Info()
		if infoErr == nil {
			fileInfo = generateFileInfo(conditions.ItemID, filePath, info)
		} else {
			fileInfo.ItemID = conditions.ItemID
			fileInfo.Name = file.Name()
			fileInfo.FilePath = filePath
			fileInfo.ParentPath = conditions.ParentPath
			fileInfo.Available = false

			if file.IsDir() {
				fileInfo.FileType = FileTypeFolder
			} else {
				fileInfo.FileType = FileTypeFile
			}
		}

		err = traverseFn(fileInfo)
		if err != nil {
			return err
		}
	}
	return err
}

func generateFileInfo(id uint, filePath string, fi fs.FileInfo) NodeFileInfo {

	var fileInfo NodeFileInfo
	fileInfo.ItemID = id
	fileInfo.Name = fi.Name()
	if fi.IsDir() {
		fileInfo.FileType = FileTypeFolder
	} else {
		fileInfo.FileType = FileTypeFile
	}
	fileInfo.ParentPath = path.Dir(filePath)
	fileInfo.FilePath = filePath
	fileInfo.Size = fi.Size()
	fileInfo.UpdatedAt = fi.ModTime()
	fileInfo.CreatedAt = fi.ModTime()
	fileInfo.Available = true

	return fileInfo
}
