package nodeitem

import (
	"errors"

	"gorm.io/gorm"

	appSample "pan/app/sample"
)

var ErrNodeItemNotFound = errors.New("nodeitem.NodeItemRepository Error: Not Found")

type NodeItemRepository interface {
	Save(NodeItem) (NodeItem, error)
	Select(uint) (NodeItem, error)
	SelectByName(string) (NodeItem, error)
	Delete(NodeItem) error
	TraverseAll(func(NodeItem) error) error
}

type nodeItemRepository struct {
	db appSample.RepositoryDB
}

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
