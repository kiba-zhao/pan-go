// Define node item repository
package nodeitem

import (
	"errors"

	"gorm.io/gorm"

	appSample "pan/app/sample"
)

var ErrNodeItemNotFound = errors.New("nodeitem.NodeItemRepository Error: Not Found")

type NodeItemRepository interface {
	// Save node item to database. If the item does not exist,
	// it will be created. If the item exists, it will be updated.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	// If the item is not found, ErrNodeItemNotFound will be returned.
	Save(NodeItem) (NodeItem, error)
	// Select a node item by its id. If the item is not found, ErrNodeItemNotFound will be returned.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	Select(uint) (NodeItem, error)
	// SelectByName retrieves a node item by its name.
	// If the item is not found, ErrNodeItemNotFound will be returned.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	SelectByName(string) (NodeItem, error)
	// Delete removes a node item from the database.
	// If the operation is successful, the item is deleted.
	// If the item is not found, ErrNodeItemNotFound will be returned.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	Delete(NodeItem) error
	// TraverseAll iterates over all NodeItems in the repository, applying the given
	// function to each item. If the database is unavailable, appSample.ErrSampleDBUnavailable
	// will be returned. If an error occurs during iteration, the iteration stops and
	// the error is returned.
	TraverseAll(func(NodeItem) error) error
	// SelectAllWithEnabled retrieves all NodeItems from the repository which have the Enabled field matching the given argument.
	// If the database is unavailable, appSample.ErrSampleDBUnavailable will be returned.
	SelectAllWithEnabled(bool) ([]NodeItem, error)
}

type nodeItemRepository struct {
	db appSample.RepositoryDB
}

// NewNodeItemRepository creates a new instance of NodeItemRepository
// using the provided appSample.RepositoryDB. This function initializes
// the repository with the given database connection, allowing operations
// to be performed on node items. If the database connection is nil,
// subsequent repository operations will return appSample.ErrSampleDBUnavailable.

func NewNodeItemRepository(db appSample.RepositoryDB) NodeItemRepository {
	return &nodeItemRepository{db: db}
}

func (repo *nodeItemRepository) Save(item NodeItem) (NodeItem, error) {
	db := repo.db
	if db == nil {
		return item, appSample.ErrSampleDBUnavailable
	}

	results := db.Save(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return item, ErrNodeItemNotFound
	}
	return item, results.Error
}

func (repo *nodeItemRepository) Select(id uint) (NodeItem, error) {
	db := repo.db
	if db == nil {
		return NodeItem{}, appSample.ErrSampleDBUnavailable
	}
	var item NodeItem
	results := db.Take(&item, id)
	if results.Error == gorm.ErrRecordNotFound {
		return item, ErrNodeItemNotFound
	}
	return item, results.Error
}

func (repo *nodeItemRepository) SelectByName(name string) (NodeItem, error) {
	db := repo.db
	if db == nil {
		return NodeItem{}, appSample.ErrSampleDBUnavailable
	}
	var item NodeItem
	results := db.Where("name = ?", name).Take(&item)
	if results.Error == gorm.ErrRecordNotFound {
		return item, ErrNodeItemNotFound
	}
	return item, results.Error
}

func (repo *nodeItemRepository) Delete(item NodeItem) error {
	db := repo.db
	if db == nil {
		return appSample.ErrSampleDBUnavailable
	}
	results := db.Delete(&item)
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrNodeItemNotFound
	}
	return results.Error
}

func (repo *nodeItemRepository) TraverseAll(traverseFn func(NodeItem) error) error {
	db := repo.db
	if db == nil {
		return appSample.ErrSampleDBUnavailable
	}
	rows, err := db.Model(&NodeItem{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var model NodeItem
		err = db.ScanRows(rows, &model)
		if err == nil {
			err = traverseFn(model)
		}
		if err != nil {
			break
		}
	}
	return err

}

func (repo *nodeItemRepository) SelectAllWithEnabled(enabled bool) ([]NodeItem, error) {
	db := repo.db
	if db == nil {
		return []NodeItem{}, appSample.ErrSampleDBUnavailable
	}
	var items []NodeItem
	results := db.Where("enabled = ?", enabled).Find(&items)
	return items, results.Error
}
