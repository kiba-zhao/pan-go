package models

type RangeSearchCondition struct {
	RangeStart int `form:"_start" binding:"omitempty"`
	RangeEnd   int `form:"_end" binding:"omitempty"`
}

type SortSearchCondition struct {
	SortField string `form:"_sort" binding:"omitempty"`
	SortOrder string `form:"_order" binding:"omitempty"`
}
