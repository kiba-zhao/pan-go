package impl

import (
	"pan/app"
	"pan/app/constant"
	"pan/extfs/errors"
	"pan/extfs/models"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TargetRepository struct {
	Provider app.RepositoryDBProvider
}

func (repo *TargetRepository) Save(target models.Target, withVersion bool) (models.Target, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return models.Target{}, constant.ErrUnavailable
	}

	if target.ID == 0 {
		result := db.Create(&target)
		return target, result.Error
	}

	var results *gorm.DB
	if !withVersion {
		results = db.Save(&target)
	} else {
		version := target.Version + 1
		results = db.Model(&target).Where("version = ?", target.Version).Updates(models.Target{
			Name:      target.Name,
			FilePath:  target.FilePath,
			Enabled:   target.Enabled,
			Version:   version,
			Available: target.Available,
		})
	}

	if results.Error == nil && results.RowsAffected != 1 {
		return target, errors.ErrNotFound
	}

	return target, results.Error
}

func (repo *TargetRepository) Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return 0, nil, constant.ErrUnavailable
	}

	if len(conditions.SortField) > 0 {
		fields := strings.Split(conditions.SortField, ",")
		orders := strings.Split(conditions.SortOrder, ",")
		for i, field := range fields {
			if len(strings.Trim(field, " ")) <= 0 {
				continue
			}
			order := false
			if len(orders) > i {
				order = strings.ToLower(orders[i]) == "desc"
			}
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: order})
		}
	}

	if len(conditions.Keyword) > 0 {
		tx := db
		keywords := strings.Split(conditions.Keyword, ",")
		for _, keyword := range keywords {
			trimKeyword := strings.Trim(keyword, " ")
			if len(trimKeyword) > 0 {
				tx = tx.Or("name like ?", "%"+keyword+"%")
				tx = tx.Or("file_path like ?", "%"+keyword+"%")
			}
		}
		db = db.Where(tx)
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

func (repo *TargetRepository) Select(id uint, version *uint8) (models.Target, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return models.Target{}, constant.ErrUnavailable
	}

	var target models.Target
	results := db.Take(&target, id)
	if results.Error == gorm.ErrRecordNotFound {
		return target, errors.ErrNotFound
	}
	if version != nil && target.Version != *version {
		return target, errors.ErrConflict
	}
	return target, results.Error
}

func (repo *TargetRepository) Delete(target models.Target) error {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}
	results := db.Where("version = ?", target.Version).Delete(&target)
	if results.Error == nil && results.RowsAffected != 1 {
		return errors.ErrNotFound
	}
	return results.Error
}
