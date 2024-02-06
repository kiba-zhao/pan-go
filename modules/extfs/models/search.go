package models

import "time"

type SearchScopeRange struct {
	Start uint `json:"start" binding:"required"`
	End   uint `json:"end" binding:"required"`
}

type SearchScope struct {
	UpdatedTime time.Time          `json:"updated_time" binding:"required"`
	Ranges      []SearchScopeRange `json:"ranges" binding:"required"`
}

type SearchCondition struct {
	Limit uint8  `form:"limit" json:"limit" binding:"omitempty"`
	Scope string `form:"scope" json:"scope" binding:"omitempty"`
}

type SearchResult[T any] struct {
	Scope string `json:"scope"`
	Items []T    `json:"items"`
}
