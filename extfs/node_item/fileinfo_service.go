// Define node file info service
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
	// TraverseWithCondition traverses all NodeFileInfo under the given condition and
	// applies the given function to each NodeFileInfo. If the given function returns
	// an error, the traversal will be stopped and the error will be returned.
	TraverseWithCondition(func(fileInfo NodeFileInfo) error, NodeFileInfoSearchCondition) error
	// Select retrieves a NodeFileInfo with the given id and file path.
	// If the file info does not exist, ErrNodeFileInfoUnavailable will be returned.
	// If the file path is invalid, ErrNodeFilePathUnavailable will be returned.
	// If the file path is not available, ErrNodeFileInfoUnavailable will be returned.
	Select(id uint, filePath string) (NodeFileInfo, error)
	// IsNotExist returns true if the given error is either a "not found" error,
	// or if the file path is unavailable.
	IsNotExist(err error) bool
}

type NodeFileInfoService struct {
	NodeFilePathService NodeFilePathInternalService
}

// IsNotExist returns true if the given error is either a "not found" error,
// or if the file path is unavailable.
func (s *NodeFileInfoService) IsNotExist(err error) bool {
	if os.IsNotExist(err) {
		return true
	}
	if errors.Is(err, ErrNodeFileInfoUnavailable) {
		return true
	}
	return s.NodeFilePathService.IsNotExist(err)
}

// Select retrieves a NodeFileInfo with the given id and file path.
// If the file info does not exist, ErrNodeFileInfoUnavailable will be returned.
// If the file path is invalid, ErrNodeFilePathUnavailable will be returned.
// If the file path is not available, ErrNodeFileInfoUnavailable will be returned.
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

// Search retrieves all NodeFileInfo under the given id and parent path.
// The function will return the total count of the results and the results
// as a slice of NodeFileInfo. If an error occurs, the error will be returned.
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

// TraverseWithCondition traverses all NodeFileInfos under the specified conditions
// and applies the provided function to each NodeFileInfo. The traversal is based on
// the item ID and parent path specified in the conditions. If the directory path
// cannot be resolved or read, an error is returned. For each file or directory found,
// a NodeFileInfo is generated and passed to the traverse function. If the traverse
// function returns an error, the traversal stops and the error is returned.

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
