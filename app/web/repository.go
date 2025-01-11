package web

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NoopDBScope(db *gorm.DB) *gorm.DB {
	return db
}

func Paginate(start int, end int) func(db *gorm.DB) *gorm.DB {

	if start <= 0 && end <= 0 {
		return NoopDBScope
	}

	return func(db *gorm.DB) *gorm.DB {
		if start > 0 {
			db = db.Offset(start)
		}

		if end > 0 {
			db = db.Limit(end - start)
		}
		return db
	}
}

func PaginateWithRangeCondition(rangeCondition *RangeCondition) func(db *gorm.DB) *gorm.DB {
	if rangeCondition == nil {
		return NoopDBScope
	}
	return Paginate(rangeCondition.RangeStart, rangeCondition.RangeEnd)
}

func OrderBy(sortField string, sortOrder string) func(db *gorm.DB) *gorm.DB {
	if len(sortField) <= 0 {
		return NoopDBScope
	}
	return func(db *gorm.DB) *gorm.DB {
		fields := strings.Split(sortField, ",")
		orders := strings.Split(sortOrder, ",")
		for i, field := range fields {
			if len(strings.Trim(field, " ")) <= 0 {
				continue
			}
			order := false
			if len(orders) > i {
				order = strings.ToLower(orders[i]) == "desc"
			}
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: order})
		}
		return db
	}
}

func OrderByWithSortCondition(sortCondition *SortCondition) func(db *gorm.DB) *gorm.DB {
	if sortCondition == nil {
		return NoopDBScope
	}
	return OrderBy(sortCondition.SortField, sortCondition.SortOrder)
}
