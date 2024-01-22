package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type RemoteFilesStateRepository interface {
	FindOne(peerId string) (models.RemoteFilesState, error)
	Save(state models.RemoteFilesState) error
}

type remoteFilesStateRepositoryImpl struct {
	DB *gorm.DB
}

func NewRemoteFilesStateRepository(db *gorm.DB) RemoteFilesStateRepository {
	repo := new(remoteFilesStateRepositoryImpl)
	repo.DB = db
	return repo
}

func (repo *remoteFilesStateRepositoryImpl) FindOne(peerId string) (models.RemoteFilesState, error) {

	var model models.RemoteFilesState
	model.ID = peerId
	results := repo.DB.Take(&model)
	if results.Error != nil {
		return model, results.Error
	}

	return model, nil
}

func (repo *remoteFilesStateRepositoryImpl) Save(state models.RemoteFilesState) error {

	results := repo.DB.Save(&state)
	return results.Error
}
