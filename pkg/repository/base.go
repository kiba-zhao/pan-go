package repository

import (
	"context"
	"sync"
)

type Repository interface {
	SetupToRepository(db RepositoryDB) error
}

type RepositoryBase struct {
	rw sync.RWMutex
	db RepositoryDB
}

var _ = (Repository)((*RepositoryBase)(nil))

func (repo *RepositoryBase) SetupToRepository(db RepositoryDB) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()
	repo.db = db
	return nil
}

func (repo *RepositoryBase) DB() RepositoryDB {
	repo.rw.RLock()
	defer repo.rw.RUnlock()
	return repo.db
}

func (repo *RepositoryBase) WithContext(ctx context.Context) RepositoryDB {
	db := repo.DB()
	if db != nil {
		return db.WithContext(ctx)
	}
	return db
}
