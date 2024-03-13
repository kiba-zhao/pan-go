package models

type Module struct {
	Avatar   string
	Name     string
	Desc     string
	Enabled  bool
	ReadOnly bool
	HasWeb   bool
}

type ModuleSearchCondition struct {
	Keyword string `form:"keyword" binding:"omitempty"`
}

// ModuleEnabled struct, Name and Enabled fields are tagged
type ModuleFields struct {
	Enabled *bool `form:"enabled" binding:"required"`
}
