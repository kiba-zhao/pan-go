package impl

import (
	"pan/app"
	"pan/app/constant"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TargetFileRepository struct {
	Provider app.RepositoryDBProvider
}

func (repo *TargetFileRepository) Save(targetFile models.TargetFile) (models.TargetFile, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return targetFile, constant.ErrUnavailable
	}

	results := db.Save(&targetFile)
	if results.Error == nil && results.RowsAffected != 1 {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) Delete(targetFile models.TargetFile) error {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}
	results := db.Delete(&targetFile)
	if results.Error == nil && results.RowsAffected != 1 {
		return errors.ErrNotFound
	}
	return results.Error
}

func (repo *TargetFileRepository) DeleteByTargetId(targetId uint) error {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}
	var targetFile models.TargetFile
	targetFile.TargetID = targetId

	results := db.Where(&targetFile).Delete(&targetFile)
	return results.Error
}

func (repo *TargetFileRepository) Select(id uint64, includeAssociated bool) (models.TargetFile, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return models.TargetFile{}, constant.ErrUnavailable
	}
	if includeAssociated {
		db = db.Preload(clause.Associations)
	}

	var targetFile models.TargetFile
	results := db.Take(&targetFile, id)
	if results.Error == gorm.ErrRecordNotFound {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) SelectByFilePathAndTargetId(filepath string, targetId uint, hashCode string, includeAssociated bool) (models.TargetFile, error) {

	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return models.TargetFile{}, constant.ErrUnavailable
	}
	if includeAssociated {
		db = db.Preload(clause.Associations)
	}

	var targetFile models.TargetFile
	targetFile.FilePath = filepath
	targetFile.TargetID = targetId
	targetFile.HashCode = hashCode

	results := db.Where(&targetFile).Take(&targetFile)
	if results.Error == gorm.ErrRecordNotFound {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) TraverseByTargetId(f repositories.TargetFileTraverse, targetId uint) error {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}

	var targetFile models.TargetFile
	targetFile.TargetID = targetId

	rows, err := db.Model(&targetFile).Where(&targetFile).Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var targetFile_ models.TargetFile
		err = db.ScanRows(rows, &targetFile_)
		if err != nil {
			break
		}
		err = f(targetFile_)
		if err != nil {
			break
		}
	}

	return err
}

func (repo *TargetFileRepository) Search(conditions models.TargetFileSearchCondition, includeAssociated bool) (total int64, items []models.TargetFile, err error) {

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
				tx = tx.Or("file_path like ?", "%"+keyword+"%")
			}
		}
		db.Where(tx)
	}

	if conditions.TargetID > 0 {
		db = db.Where("target_id = ?", conditions.TargetID)
	}

	results := db.Model(&models.TargetFile{}).Count(&total)
	if results.Error != nil {
		return
	}

	if includeAssociated {
		db = db.Preload(clause.Associations)
	}

	if conditions.RangeStart > 0 {
		db = db.Limit(conditions.RangeStart)
	}

	if conditions.RangeEnd > 0 {
		db = db.Limit(conditions.RangeEnd - conditions.RangeStart)
	}

	results = db.Find(&items)
	err = results.Error
	return
}
