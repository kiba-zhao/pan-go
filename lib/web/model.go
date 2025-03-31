package web

// Range Condition for query
type RangeCondition struct {
	RangeStart int `form:"_start" binding:"omitempty"`
	RangeEnd   int `form:"_end" binding:"omitempty"`
}

// Sort Condition for query
type SortCondition struct {
	SortField string `form:"_sort" binding:"omitempty"`
	SortOrder string `form:"_order" binding:"omitempty"`
}
