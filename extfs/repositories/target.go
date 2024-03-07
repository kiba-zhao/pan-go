package repositories

import (
	"pan/extfs/models"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TargetRepository interface {
	Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error)
}

type targetRepositoryImpl struct {
	DB *gorm.DB
}

func NewTargetRepository(db *gorm.DB) TargetRepository {
	return &targetRepositoryImpl{DB: db}
}

func (repo *targetRepositoryImpl) Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error) {
	db := repo.DB

	if len(conditions.SortField) > 0 {
		db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: conditions.SortField}, Desc: conditions.SortOrder})
	}

	if len(conditions.Keyword) > 0 {
		tx := repo.DB
		keywords := strings.Split(conditions.Keyword, ",")
		for _, keyword := range keywords {
			trimKeyword := strings.Trim(keyword, " ")
			if len(trimKeyword) > 0 {
				tx = tx.Or("name like ?", "%"+keyword+"%")
				tx = tx.Or("file_path like ?", "%"+keyword+"%")
			}
		}
		db.Where(tx)
	}

	if conditions.Enabled != nil {
		db = db.Where("enabled = ?", *conditions.Enabled)
	}

	results := db.Model(&models.Target{}).Count(&total)
	if results.Error != nil {
		return
	}

	if conditions.RangeStart > 0 {
		db = db.Offset(conditions.RangeStart)
	}

	if conditions.RangeEnd > 0 {
		db = db.Limit(conditions.RangeEnd - conditions.RangeStart)
	}

	results = db.Find(&items)
	err = results.Error
	return
}
