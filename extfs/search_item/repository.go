package searchitem

import (
	"errors"
	appSample "pan/app/sample"

	"gorm.io/gorm"
)

var ErrSearchItemNotFound = errors.New("searchitem.SearchItemRepository Error: Not Found")

type SearchItemRepository interface {
	Save(SearchItem) (SearchItem, error)
	Select(id uint) (SearchItem, error)
	Delete(id uint) error
	Search(SearchItemCondition) (int64, []SearchItem, error)
}

type searchItemRepositoryImpl struct {
	db appSample.RepositoryDB
}

func NewSearchItemRepository(db appSample.RepositoryDB) SearchItemRepository {
	return &searchItemRepositoryImpl{db: db}
}

func (repo *searchItemRepositoryImpl) Save(item SearchItem) (SearchItem, error) {
	db := repo.db
	if db == nil {
		return item, appSample.ErrSampleDBUnavailable
	}
	results := db.Save(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return item, ErrSearchItemNotFound
	}
	return item, results.Error
}

func (repo *searchItemRepositoryImpl) Select(id uint) (SearchItem, error) {
	db := repo.db
	if db == nil {
		return SearchItem{}, appSample.ErrSampleDBUnavailable
	}
	var item SearchItem
	results := db.Take(&item, id)
	if results.Error == gorm.ErrRecordNotFound {
		return item, ErrSearchItemNotFound
	}
	return item, results.Error
}

func (repo *searchItemRepositoryImpl) Delete(id uint) error {
	db := repo.db
	if db == nil {
		return appSample.ErrSampleDBUnavailable
	}
	results := db.Delete(&SearchItem{ID: id})
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrSearchItemNotFound
	}
	return results.Error
}

func (repo *searchItemRepositoryImpl) Search(condition SearchItemCondition) (int64, []SearchItem, error) {
	db := repo.db
	if db == nil {
		return 0, nil, appSample.ErrSampleDBUnavailable
	}
	if len(condition.Query) > 0 {
		db = db.Where("query like ?", condition.Query+"%")
	}

	total := int64(0)
	results := db.Model(&SearchItem{}).Count(&total)

	if results.Error != nil || total <= 0 {
		return total, nil, results.Error
	}

	if condition.RangeStart > 0 {
		db = db.Offset(condition.RangeStart)
	}

	if condition.RangeEnd > 0 {
		db = db.Limit(condition.RangeEnd - condition.RangeStart)
	}

	var items []SearchItem
	results = db.Order("updated_at desc").Find(&items)
	return total, items, results.Error
}
