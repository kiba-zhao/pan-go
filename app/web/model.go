package web

type RangeCondition struct {
	RangeStart int `form:"_start" binding:"omitempty"`
	RangeEnd   int `form:"_end" binding:"omitempty"`
}

type SortCondition struct {
	SortField string `form:"_sort" binding:"omitempty"`
	SortOrder string `form:"_order" binding:"omitempty"`
}
