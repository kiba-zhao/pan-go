package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type TargetRepository interface {
	FindAllWithEnabled() (targets []models.Target, err error)
	Save(target models.Target) error
}

type targetRepositoryImpl struct {
	DB *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	repo := new(targetRepositoryImpl)
	repo.DB = db
	return repo
}

func (repo *targetRepositoryImpl) FindAllWithEnabled() (targets []models.Target, err error) {

	results := repo.DB.Where("enabled = ?", true).Find(&targets)
	return targets, results.Error
}

func (repo *targetRepositoryImpl) Save(target models.Target) error {
	results := repo.DB.Save(&target)
	return results.Error
}
