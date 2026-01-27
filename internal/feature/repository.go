package feature

import (
	"errors"
	"pan/internal/repository"
	"sync"
)

var ErrFeatureRepositoryDBUnavailable = errors.New("feature.Repository Error: Database unavailable")

type Repository struct {
	rw sync.RWMutex
	db repository.RepositoryDB
}

var _ = (repository.Repository)((*Repository)(nil))

func (repo *Repository) SetupToRepository(db repository.RepositoryDB) error {
	repo.rw.Lock()
	defer repo.rw.Unlock()
	repo.db = db
	return nil
}

func (repo *Repository) DB() repository.RepositoryDB {
	repo.rw.RLock()
	defer repo.rw.RUnlock()
	return repo.db
}
