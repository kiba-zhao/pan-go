package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type FilesStateRepository interface {
	GetLastOne() (models.FilesState, error)
}
type filesStateRepositoryImpl struct {
	DB *gorm.DB
}

func NewFilesStateRepository(db *gorm.DB) FilesStateRepository {
	repo := new(filesStateRepositoryImpl)
	repo.DB = db
	return repo
}

func (repo *filesStateRepositoryImpl) GetLastOne() (models.FilesState, error) {
	var model models.FilesState
	results := repo.DB.Last(&model)
	if results.Error != nil {
		return model, results.Error
	}
	return model, nil
}
