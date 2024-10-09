package models

import "time"

const (
	FILETYPE_FOLDER = "D"
	FILETYPE_FILE   = "F"
)

type DiskFile struct {
	ID         string    `json:"id" form:"id"`
	Name       string    `json:"name" form:"name"`
	FilePath   string    `json:"filePath" form:"filePath"`
	ParentPath string    `json:"parentPath" form:"parentPath"`
	FileType   string    `json:"fileType" form:"fileType"`
	UpdatedAt  time.Time `json:"updatedAt" form:"updatedAt"`
}

type DiskFileSearchCondition struct {
	ParentPath string `form:"parentPath" binding:"omitempty"`
	FilePath   string `form:"filePath" binding:"omitempty"`
	FileType   string `form:"fileType" binding:"omitempty"`
}
