package impl

import (
	"pan/app"
	"pan/app/constant"
	"pan/extfs/models"

	"gorm.io/gorm"
)

type NodeItemRepository struct {
	Provider app.RepositoryDBProvider
}

func (repo *NodeItemRepository) Save(item models.NodeItem) (models.NodeItem, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return item, constant.ErrUnavailable
	}

	results := db.Save(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return item, constant.ErrNotFound
	}
	return item, results.Error
}

func (repo *NodeItemRepository) Select(id uint) (models.NodeItem, error) {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return models.NodeItem{}, constant.ErrUnavailable
	}
	var item models.NodeItem
	results := db.Take(&item, id)
	if results.Error == gorm.ErrRecordNotFound {
		return item, constant.ErrNotFound
	}
	return item, results.Error
}

func (repo *NodeItemRepository) Delete(item models.NodeItem) error {
	db := app.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}
	results := db.Delete(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return constant.ErrNotFound
	}
	return results.Error
}
