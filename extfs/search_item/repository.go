package searchitem

import (
	"errors"
	appSample "pan/app/sample"
	"pan/app/web"

	"gorm.io/gorm"
)

var ErrSearchItemNotFound = errors.New("searchitem.SearchItemRepository Error: Not Found")

type SearchItemRepository interface {
	Save(SearchItem) (SearchItem, error)
	Create(SearchItem) (SearchItem, error)
	SelectOrCreate(SearchItem) (SearchItem, bool, error)
	Select(id uint64) (SearchItem, error)
	Delete(id uint64) error
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

func (repo *searchItemRepositoryImpl) Create(fields SearchItem) (SearchItem, error) {
	db := repo.db
	if db == nil {
		return SearchItem{}, appSample.ErrSampleDBUnavailable
	}
	var item SearchItem
	results := db.Create(&item)
	return item, results.Error
}

func (repo *searchItemRepositoryImpl) SelectOrCreate(item SearchItem) (SearchItem, bool, error) {
	db := repo.db
	if db == nil {
		return item, false, appSample.ErrSampleDBUnavailable
	}
	db = db.Where("query = ?", item.Query)
	db = db.Order("updated_at desc")
	results := db.FirstOrCreate(&item)
	if results.Error == nil {
		return item, results.RowsAffected == 1, results.Error
	}
	return item, false, results.Error
}

func (repo *searchItemRepositoryImpl) Select(id uint64) (SearchItem, error) {
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

func (repo *searchItemRepositoryImpl) Delete(id uint64) error {
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
	if len(condition.Q) > 0 {
		db = db.Where("query like ?", condition.Q+"%")
	}
	if len(condition.Query) > 0 {
		db = db.Where("query = ?", condition.Query)
	}

	total := int64(0)
	results := db.Model(&SearchItem{}).Count(&total)

	if results.Error != nil || total <= 0 {
		return total, nil, results.Error
	}

	db = db.Scopes(web.PaginateWithRangeCondition(&condition.RangeCondition))

	var items []SearchItem
	results = db.Order("updated_at desc").Find(&items)
	return total, items, results.Error
}
