package searchitem

import (
	"errors"
	"pan/lib/feature"
	"pan/lib/repository"
	"pan/lib/web"

	"gorm.io/gorm"
)

var ErrSearchItemNotFound = errors.New("searchitem.SearchItemRepository Error: Not Found")

type SearchItemRepository interface {
	repository.Repository
	// Save saves the given SearchItem to the database.
	// It returns the saved SearchItem and an error if any issues occur during the query.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	Save(SearchItem) (SearchItem, error)
	// Create creates the given SearchItem in the database.
	// It returns the created SearchItem and an error if any issues occur during the query.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	Create(SearchItem) (SearchItem, error)
	// SelectOrCreate returns the SearchItem with the given id.
	// If the SearchItem does not exist, it creates a new one.
	// It returns the SearchItem, a boolean indicating whether the SearchItem was created, and an error if any issues occur during the query.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	SelectOrCreate(SearchItem) (SearchItem, bool, error)
	// Select retrieves the SearchItem associated with the given id.
	// It returns the retrieved SearchItem and an error if any issues occur during the query.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	// If the SearchItem was not found, it returns ErrSearchItemNotFound.
	//
	Select(id uint64) (SearchItem, error)
	// Delete removes the SearchItem associated with the given id from the database.
	// It returns an error if any issues occur during the deletion process.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	// If the SearchItem was not found, it returns ErrSearchItemNotFound.
	//
	Delete(id uint64) error
	// Search retrieves SearchItems associated with the given condition, applying the given range condition for pagination.
	// It returns the total number of results, a slice of SearchItems, and an error if any issues occur during the query.
	// If the database is unavailable, it returns feature.ErrFeatureRepositoryDBUnavailable.
	//
	Search(SearchItemCondition) (int64, []SearchItem, error)
}

type stdSearchItemRepository struct {
	feature.Repository
}

func NewSearchItemRepository() SearchItemRepository {
	return &stdSearchItemRepository{}
}

func (repo *stdSearchItemRepository) Save(item SearchItem) (SearchItem, error) {
	db := repo.DB()
	if db == nil {
		return item, feature.ErrFeatureRepositoryDBUnavailable
	}
	results := db.Save(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return item, ErrSearchItemNotFound
	}
	return item, results.Error
}

func (repo *stdSearchItemRepository) Create(item SearchItem) (SearchItem, error) {
	db := repo.DB()
	if db == nil {
		return SearchItem{}, feature.ErrFeatureRepositoryDBUnavailable
	}
	results := db.Create(&item)
	return item, results.Error
}

func (repo *stdSearchItemRepository) SelectOrCreate(item SearchItem) (SearchItem, bool, error) {
	db := repo.DB()
	if db == nil {
		return item, false, feature.ErrFeatureRepositoryDBUnavailable
	}
	db = db.Where("query = ?", item.Query)
	db = db.Order("updated_at desc")
	results := db.FirstOrCreate(&item)
	if results.Error == nil {
		return item, results.RowsAffected == 1, results.Error
	}
	return item, false, results.Error
}

func (repo *stdSearchItemRepository) Select(id uint64) (SearchItem, error) {
	db := repo.DB()
	if db == nil {
		return SearchItem{}, feature.ErrFeatureRepositoryDBUnavailable
	}
	var item SearchItem
	results := db.Take(&item, id)
	if results.Error == gorm.ErrRecordNotFound {
		return item, ErrSearchItemNotFound
	}
	return item, results.Error
}

func (repo *stdSearchItemRepository) Delete(id uint64) error {
	db := repo.DB()
	if db == nil {
		return feature.ErrFeatureRepositoryDBUnavailable
	}
	results := db.Delete(&SearchItem{ID: id})
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrSearchItemNotFound
	}
	return results.Error
}

func (repo *stdSearchItemRepository) Search(condition SearchItemCondition) (int64, []SearchItem, error) {
	db := repo.DB()
	if db == nil {
		return 0, nil, feature.ErrFeatureRepositoryDBUnavailable
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
