package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type ExtFSRepository interface {
	GetLatestOne() (models.ExtFS, error)
}

type extFSRepositoryImpl struct {
	DB *gorm.DB
}

// NewExtFSRepository ...
func NewExtFSRepository(db *gorm.DB) ExtFSRepository {
	repo := new(extFSRepositoryImpl)
	repo.DB = db
	return repo
}

// GetLatestOne ...
func (repo *extFSRepositoryImpl) GetLatestOne() (models.ExtFS, error) {
	var extFS models.ExtFS
	results := repo.DB.Last(&extFS)
	if results.Error != nil {
		return extFS, results.Error
	}
	return extFS, nil
}
