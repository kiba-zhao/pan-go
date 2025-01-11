package appnode

import (
	"errors"
	"pan/app/sample"
	"pan/app/web"
	"strings"

	"gorm.io/gorm"
)

var ErrAppNodeNotFound = errors.New("appnode.AppNodeRepository Error: Not Found")

type AppNodeRepository interface {
	Search(AppNodeSearchCondition) (int64, []AppNode, error)
	Save(AppNode) (AppNode, error)
	Select(uint) (AppNode, error)
	SelectByName(string) (AppNode, error)
	Delete(AppNode) error
	SelectByPeerID(string) (AppNode, error)
	TraverseAll(func(AppNode) error) error
}

type peerNodeRepository struct {
	db *gorm.DB
}

func NewAppNodeRepository(db *gorm.DB) AppNodeRepository {
	return &peerNodeRepository{db: db}
}

func (repo *peerNodeRepository) Search(conditions AppNodeSearchCondition) (int64, []AppNode, error) {

	db := repo.db
	if db == nil {
		return 0, nil, sample.ErrSampleDBUnavailable
	}

	if len(conditions.Keyword) > 0 {
		tx := db
		keywords := strings.Split(conditions.Keyword, ",")
		for _, keyword := range keywords {
			trimKeyword := strings.Trim(keyword, " ")
			if len(trimKeyword) > 0 {
				tx = tx.Or("name like ?", "%"+keyword+"%")
				tx = tx.Or("peerId = ?", keyword)
			}
		}
		db = db.Where(tx)
	}

	if conditions.Blocked != nil {
		db = db.Where("blocked = ?", *conditions.Blocked)
	}

	total := int64(0)
	results := db.Model(&AppNode{}).Count(&total)

	if results.Error != nil || total <= 0 {
		return total, nil, results.Error
	}

	db = db.Scopes(web.OrderByWithSortCondition(&conditions.SortCondition))
	db = db.Scopes(web.PaginateWithRangeCondition(&conditions.RangeCondition))

	var items []AppNode
	results = db.Find(&items)
	return total, items, results.Error

}

func (repo *peerNodeRepository) Save(node AppNode) (AppNode, error) {
	db := repo.db
	if db == nil {
		return node, sample.ErrSampleDBUnavailable
	}

	results := db.Save(&node)
	if results.Error == nil && results.RowsAffected != 1 {
		return node, ErrAppNodeNotFound
	}
	return node, results.Error
}

func (repo *peerNodeRepository) Select(id uint) (AppNode, error) {

	db := repo.db
	if db == nil {
		return AppNode{}, sample.ErrSampleDBUnavailable
	}
	var peerNode AppNode
	results := db.Take(&peerNode, id)
	if results.Error == gorm.ErrRecordNotFound {
		return peerNode, ErrAppNodeNotFound
	}
	return peerNode, results.Error
}

func (repo *peerNodeRepository) SelectByName(name string) (AppNode, error) {

	db := repo.db
	if db == nil {
		return AppNode{}, sample.ErrSampleDBUnavailable
	}
	var peerNode AppNode
	results := db.Where("name = ?", name).Take(&peerNode)
	if results.Error == gorm.ErrRecordNotFound {
		return peerNode, ErrAppNodeNotFound
	}
	return peerNode, results.Error
}

func (repo *peerNodeRepository) Delete(node AppNode) error {

	db := repo.db
	if db == nil {
		return sample.ErrSampleDBUnavailable
	}
	results := db.Delete(&node)
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrAppNodeNotFound
	}
	return results.Error
}

func (repo *peerNodeRepository) SelectByPeerID(peerId string) (AppNode, error) {

	db := repo.db
	if db == nil {
		return AppNode{}, sample.ErrSampleDBUnavailable
	}
	var peerNode AppNode
	peerNode.PeerID = peerId
	results := db.Where(&peerNode).Take(&peerNode)
	if results.Error == gorm.ErrRecordNotFound {
		return peerNode, ErrAppNodeNotFound
	}
	return peerNode, results.Error
}

func (repo *peerNodeRepository) TraverseAll(traverse func(AppNode) error) error {

	db := repo.db
	if db == nil {
		return sample.ErrSampleDBUnavailable
	}

	rows, err := db.Model(&AppNode{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var model AppNode
		err = db.ScanRows(rows, &model)
		if err == nil {
			err = traverse(model)
		}
		if err != nil {
			break
		}
	}
	return err

}
