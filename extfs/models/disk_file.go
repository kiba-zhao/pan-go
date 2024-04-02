package models

import "time"

const (
	FILETYPE_FOLDER = "D"
	FILETYPE_FILE   = "F"
)

type DiskFile struct {
	ID        string    `json:"id" form:"id"`
	Name      string    `json:"name" form:"name"`
	FilePath  string    `json:"filepath" form:"filepath"`
	Parent    string    `json:"parent" form:"parent"`
	FileType  string    `json:"fileType" form:"fileType"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
}

type DiskFileSearchCondition struct {
	Parent   string `form:"parent" binding:"omitempty"`
	FilePath string `form:"filepath" binding:"omitempty"`
}
