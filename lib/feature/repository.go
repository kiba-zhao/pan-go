package feature

import (
	"errors"
	"pan/lib/repository"
	"sync"
)

var ErrFeatureRepositoryDBUnavailable = errors.New("feature.RepositoryDB Error: Unavailable")

type Repository struct {
	db   repository.RepositoryDB
	dbRW sync.RWMutex
}

var _ = (repository.Repository)((*Repository)(nil))

func (repo *Repository) SetupToRepository(db repository.RepositoryDB) error {
	repo.dbRW.Lock()
	defer repo.dbRW.Unlock()
	repo.db = db
	return nil
}

func (repo *Repository) DB() repository.RepositoryDB {
	repo.dbRW.RLock()
	defer repo.dbRW.RUnlock()
	return repo.db
}

type RepositoryMeta = Metadata[repository.Repository]

func NewRepositoryMeta[T any](repo repository.Repository) RepositoryMeta {
	return newMetadata[T, repository.Repository](repo)
}

type RepositoryMetaProvider interface {
	RepositoryMetaList() []RepositoryMeta
}

func getRepositoryMetaList(feature interface{}) []RepositoryMeta {
	if metaProvider, ok := feature.(RepositoryMetaProvider); ok {
		return metaProvider.RepositoryMetaList()
	}
	return nil
}

func getRepository(feature interface{}) repository.Repository {
	if repository, ok := feature.(repository.Repository); ok {
		return repository
	}
	return nil
}
