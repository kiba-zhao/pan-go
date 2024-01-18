package repositories

import (
	"pan/modules/extfs/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository interface {
	Find(condition *models.TagFindCondition) ([]models.Tag, error)
}

type tagRepository struct {
	DB *gorm.DB
}

// New ...
func NewTagRepository(db *gorm.DB) TagRepository {
	repo := new(tagRepository)
	repo.DB = db
	return repo
}

// Find ...
func (r *tagRepository) Find(condition *models.TagFindCondition) (tags []models.Tag, err error) {
	db := r.DB
	if condition != nil {

		owner := uuid.Nil
		if len(condition.Owner) > 0 {
			owner, _ = uuid.Parse(condition.Owner)
		}
		if len(condition.Name) > 0 && owner != uuid.Nil {
			db = db.Where("owner = ? And name like ?", condition.Owner, condition.Name+"%")
		} else if owner != uuid.Nil {
			db = db.Where("owner = ?", condition.Owner)
		} else if len(condition.Name) > 0 {
			db = db.Where("name like ?", condition.Name+"%")
		}

		if condition.Limit > 0 {
			db = db.Limit(condition.Limit)
		}
		if condition.Offset > 0 {
			db = db.Offset(condition.Offset)
		}
	}

	results := db.Find(&tags)
	if results.Error != nil {
		err = results.Error
	}

	return
}
