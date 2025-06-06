package feature

import (
	"errors"
	"pan/lib/repository"
	"reflect"
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

type RepositoryMeta struct {
	Type       reflect.Type
	Repository repository.Repository
}

func NewRepositoryMeta[T any](repo repository.Repository) RepositoryMeta {
	return RepositoryMeta{
		Type:       reflect.TypeFor[T](),
		Repository: repo,
	}
}

type RepositoryMetaProvider interface {
	RepositoryMetaList() []RepositoryMeta
}

var _ = (repository.RepositoryClusterModule)((*stdFeatureModule)(nil))

func (module *stdFeatureModule) DBName() string {
	featureHelper := module.featureHelper
	if repositoryClusterModule, ok := featureHelper.feature.(repository.RepositoryClusterModule); ok {
		return repositoryClusterModule.DBName()
	}
	return featureHelper.FeatureName() + ".db"
}

func (module *stdFeatureModule) SetupToRepository(db repository.RepositoryDB) error {
	featureHelper := module.featureHelper
	if repository := getRepository(featureHelper.feature); repository != nil {
		err := repository.SetupToRepository(db)
		if err != nil {
			return err
		}
	}

	metaList := getRepositoryMetaList(featureHelper.feature)
	if len(metaList) <= 0 {
		return nil
	}

	for _, meta := range metaList {
		err := meta.Repository.SetupToRepository(db)
		if err != nil {
			return err
		}
	}
	return nil
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
