package models

type RangeSearchCondition struct {
	RangeStart int `form:"range-start" binding:"omitempty"`
	RangeEnd   int `form:"range-end" binding:"omitempty"`
}

type SortSearchCondition struct {
	SortField string `form:"sort-field" binding:"omitempty"`
	SortOrder bool   `form:"sort-order" binding:"omitempty"`
}
