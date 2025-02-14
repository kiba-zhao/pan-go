package nodeitem

import (
	"net/http"
	"os"
)

func GenerateMimeTypeWithFilePath(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return GenerateMimeType(file)
}

func GenerateMimeType(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	fileType := http.DetectContentType(buffer)
	return fileType, nil
}
