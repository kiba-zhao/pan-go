package repositories

import "gorm.io/gorm"

type RemoteFileRepository interface {
}

type remoteFileRepositoryImpl struct {
	DB *gorm.DB
}

func NewRemoteFileRepository(db *gorm.DB) RemoteFileRepository {
	repo := new(remoteFileRepositoryImpl)
	repo.DB = db
	return repo
}
