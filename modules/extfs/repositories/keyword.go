package repositories

import (
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type KeywordRepository interface {
	Find(condition *models.KeywordFindCondition) ([]models.Keyword, error)
}

type keywordRepository struct {
	DB *gorm.DB
}

// New ...
func NewKeywordRepository(db *gorm.DB) KeywordRepository {
	repo := new(keywordRepository)
	repo.DB = db
	return repo
}

// Find ...
func (r *keywordRepository) Find(condition *models.KeywordFindCondition) (keywords []models.Keyword, err error) {
	db := r.DB
	if condition != nil {

		if len(condition.Keyword) > 0 {
			db = db.Where("name like ?", condition.Keyword+"%")
		}

		if condition.Limit > 0 {
			db = db.Limit(condition.Limit)
		}
		if condition.Offset > 0 {
			db = db.Offset(condition.Offset)
		}
	}

	results := db.Find(&keywords)
	if results.Error != nil {
		err = results.Error
	}

	return
}
