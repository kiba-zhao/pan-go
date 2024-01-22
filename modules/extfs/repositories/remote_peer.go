package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type RemotePeerRepository interface {
	FindOne(id string) (models.RemotePeer, error)
}

type remotePeerRepositoryImpl struct {
	DB *gorm.DB
}

// NewRemotePeerRepository ...
func NewRemotePeerRepository(db *gorm.DB) RemotePeerRepository {
	repo := new(remotePeerRepositoryImpl)
	repo.DB = db
	return repo
}

// FineOne ...
func (repo *remotePeerRepositoryImpl) FindOne(id string) (models.RemotePeer, error) {

	var model models.RemotePeer
	model.ID = id
	results := repo.DB.Take(&model)
	if results.Error != nil {
		return model, results.Error
	}

	return model, nil
}
