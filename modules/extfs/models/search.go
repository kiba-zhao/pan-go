package models

type SearchCondition struct {
	Offset int `form:"offset" json:"offset" binding:"omitempty"`
	Limit  int `form:"limit" json:"limit" binding:"omitempty"`
}
type SearchResult[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}
