package repository

import (
	"errors"

	"gorm.io/gorm"
)

type RepositoryDB = *gorm.DB

var ErrRepositoryDBUnavailable = errors.New("repository.RepositoryDB Error: Unavailable")
var ErrRepositoryRecordNotFound = gorm.ErrRecordNotFound
