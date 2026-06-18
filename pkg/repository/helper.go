package repository

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NoopDBScope is a noop function that returns the passed-in *gorm.DB without modifying it.
// This is useful as a placeholder when a function is expected to return a scope, but no
// actual scoping is needed.
func NoopDBScope(db *gorm.DB) *gorm.DB {
	return db
}

// Paginate returns a scope function that paginates the passed-in *gorm.DB by setting the Offset and Limit clauses.
// If start is 0 and end is 0, Paginate returns a noop scope function.
// If start is greater than 0, Paginate sets the Offset clause to start.
// If end is greater than 0, Paginate sets the Limit clause to end - start.
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

// OrderBy returns a scope function that sets the ORDER BY clause on the passed-in *gorm.DB.
// If sortField is empty, OrderBy returns a noop scope function.
// If sortField is not empty, OrderBy splits sortField by comma and uses each field as a column to order by.
// If sortOrder is not empty, OrderBy splits sortOrder by comma and uses each value as the order for the corresponding field in sortField.
// If a field in sortField does not have a corresponding value in sortOrder, the order for that field defaults to ascending.
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

func TableName(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table(tableName)
	}
}
