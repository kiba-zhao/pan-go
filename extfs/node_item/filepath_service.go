// Define node file path service
package nodeitem

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrNodeFilePathUnavailable = errors.New("nodeitem.NodeFilePathService Error: Unavailable")
var ErrNodeFilePathWithoutFolder = errors.New("nodeitem.NodeFilePathService Error: Without Folder")

type NodeFilePathInternalService interface {
	// IsNotExist returns true if the given error is either a "not found" error,
	// or if the file path is unavailable.
	IsNotExist(err error) bool
	// Select retrieves the real file path for a given node item ID and file path.
	// If the node item is not available, ErrNodeFilePathUnavailable will be returned.
	// If the file path is not specified, the node item's root file path is returned.
	// Returns the resolved file path and an error if any issues occur during retrieval.
	Select(id uint, filePath string) (string, error)
	// SelectWithoutFolder retrieves the real file path for a given node item ID and file path,
	// without traversing into the folder if the file path is a folder.
	// If the node item is not available, ErrNodeFilePathUnavailable will be returned.
	// If the file path is not specified, the node item's root file path is returned.
	// If the file path is a folder, ErrNodeFilePathWithoutFolder will be returned.
	// Returns the resolved file path and an error if any issues occur during retrieval.
	SelectWithoutFolder(id uint, filePath string) (string, error)
}

type NodeFilePathService struct {
	NodeItemService NodeItemInternalService
}

// IsNotExist returns true if the given error is either a "not found" error,
// or if the file path is unavailable.
func (s *NodeFilePathService) IsNotExist(err error) bool {
	return s.NodeItemService.IsNotExist(err)
}

// Select retrieves the real file path for a given node item ID and file path.
// If the node item is not available, ErrNodeFilePathUnavailable will be returned.
// If the file path is not specified, the node item's root file path is returned.
// Returns the resolved file path and an error if any issues occur during retrieval.
func (s *NodeFilePathService) Select(id uint, filePath string) (string, error) {
	nodeItem, err := s.NodeItemService.Select(id)
	if err != nil {
		return "", err
	}
	if !nodeItem.Available {
		return "", ErrNodeFilePathUnavailable
	}

	if len(filePath) <= 0 {
		return nodeItem.FilePath, nil
	}

	filePath_ := filepath.FromSlash(filePath)
	realFilePath := filepath.Join(nodeItem.FilePath, filePath_)
	return realFilePath, err

}

// SelectWithoutFolder retrieves the real file path for a given node item ID and file path.
// If the node item is not available, ErrNodeFilePathUnavailable will be returned.
// If the file path is not specified, the node item's root file path is returned.
// If the file path is a folder, ErrNodeFilePathWithoutFolder will be returned.
// Returns the resolved file path and an error if any issues occur during retrieval.

func (s *NodeFilePathService) SelectWithoutFolder(id uint, filePath string) (string, error) {
	realFilePath, err := s.Select(id, filePath)
	if err != nil {
		return "", err
	}
	fileStat, err := os.Stat(realFilePath)
	if err != nil {
		return realFilePath, err
	}
	if fileStat.IsDir() {
		return "", ErrNodeFilePathWithoutFolder
	}
	return realFilePath, err
}
