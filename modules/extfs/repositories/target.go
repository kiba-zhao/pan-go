package repositories

import (
	"pan/modules/extfs/models"
	"strings"

	"gorm.io/gorm"
)

type TargetRepository interface {
	FindAllWithEnabled() (targets []models.Target, err error)
	Save(target models.Target) error
	Search(condition *models.TargetSearchCondition) (int64, []models.Target, error)
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

	results := repo.DB.Find(&targets, models.Target{Enabled: true})
	return targets, results.Error
}

func (repo *targetRepositoryImpl) Save(target models.Target) error {
	results := repo.DB.Save(&target)
	return results.Error
}

func (repo *targetRepositoryImpl) Search(condition *models.TargetSearchCondition) (total int64, targets []models.Target, err error) {

	db := repo.DB
	db = db.Order("modify_time desc")

	kw := strings.Trim(condition.Keyword, " ")
	if len(kw) > 0 {
		tx := repo.DB
		keywords := strings.Split(kw, ",")
		for _, keyword := range keywords {
			trimKeyword := strings.Trim(keyword, " ")
			if len(trimKeyword) > 0 {
				tx = tx.Or("name like ?", "%"+keyword+"%")
				tx = tx.Or("file_path like ?", "%"+keyword+"%")
			}
		}
		db.Where(tx)
	}

	results := db.Model(&models.Target{}).Count(&total)
	if results.Error != nil {
		err = results.Error
		return
	}

	db = db.Limit(condition.Limit)
	db = db.Offset(condition.Offset)

	results = db.Find(&targets)
	err = results.Error
	return
}
