// Define mime type utility for node item
package nodeitem

import (
	"net/http"
	"os"
)

// GenerateMimeTypeWithFilePath generates the mime type of the file by reading the first 512 bytes
// from the file at the given filePath.
func GenerateMimeTypeWithFilePath(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return GenerateMimeType(file)
}

// GenerateMimeType reads the first 512 bytes from file and return the mime type
// according to the content. If the file is not readable, it returns an error.
func GenerateMimeType(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)
	return fileType, nil
}
