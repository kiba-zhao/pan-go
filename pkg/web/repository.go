package web

import (
	"pan/pkg/repository"

	"gorm.io/gorm"
)

// PaginateWithRangeCondition returns a scope function that paginates the passed-in *gorm.DB with the passed-in RangeCondition.
// If RangeCondition is nil, PaginateWithRangeCondition returns a noop scope function.
// If RangeCondition is not nil, PaginateWithRangeCondition calls Paginate with RangeCondition's RangeStart and RangeEnd fields.
func PaginateWithRangeCondition(rangeCondition *RangeCondition) func(db *gorm.DB) *gorm.DB {
	if rangeCondition == nil {
		return repository.NoopDBScope
	}
	return repository.Paginate(rangeCondition.RangeStart, rangeCondition.RangeEnd)
}

// OrderByWithSortCondition returns a scope function that sets the ORDER BY clause on the passed-in *gorm.DB with the passed-in SortCondition.
// If SortCondition is nil, OrderByWithSortCondition returns a noop scope function.
// If SortCondition is not nil, OrderByWithSortCondition calls OrderBy with SortCondition's SortField and SortOrder fields.
func OrderByWithSortCondition(sortCondition *SortCondition) func(db *gorm.DB) *gorm.DB {
	if sortCondition == nil {
		return repository.NoopDBScope
	}
	return repository.OrderBy(sortCondition.SortField, sortCondition.SortOrder)
}
