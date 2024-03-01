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

type ModuleSearchResult struct {
	Total int
	Items []Module
}