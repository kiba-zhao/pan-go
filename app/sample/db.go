package sample

import (
	"errors"

	"gorm.io/gorm"
)

// RepositoryDB is database type used by sample module
type RepositoryDB = *gorm.DB

var ErrSampleDBUnavailable = errors.New("sample.RepositoryDB Error: Unavailable")
