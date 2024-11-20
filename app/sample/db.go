package sample

import (
	"errors"

	"gorm.io/gorm"
)

type RepositoryDB = *gorm.DB

type RepositoryDBProvider interface {
	DB() RepositoryDB
}

func DBForProvider(provider RepositoryDBProvider) RepositoryDB {
	if provider == nil {
		return nil
	}
	return provider.DB()
}

var ErrSampleDBUnavailable = errors.New("sample.RepositoryDB Error: Unavailable")
