package sample

import (
	"errors"

	"gorm.io/gorm"
)

type RepositoryDB = *gorm.DB

var ErrSampleDBUnavailable = errors.New("sample.RepositoryDB Error: Unavailable")
