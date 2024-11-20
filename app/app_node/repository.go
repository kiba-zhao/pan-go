package appnode

import (
	"errors"
	"pan/app/sample"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAppNodeNotFound = errors.New("appnode.AppNodeRepository Error: Not Found")

type AppNodeRepository interface {
	Search(AppNodeSearchCondition) (int64, []AppNode, error)
	Save(AppNode) (AppNode, error)
	Select(uint) (AppNode, error)
	SelectByName(string) (AppNode, error)
	Delete(AppNode) error
	SelectByPeerID(string) (AppNode, error)
	TraverseWithPeerIDs(func(AppNode) error, []string) error
}

type peerNodeRepository struct {
	Provider sample.RepositoryDBProvider
}

func NewAppNodeRepository(provider sample.RepositoryDBProvider) AppNodeRepository {
	return &peerNodeRepository{Provider: provider}
}

func (repo *peerNodeRepository) Search(conditions AppNodeSearchCondition) (int64, []AppNode, error) {

	db := sample.DBForProvider(repo.Provider)
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

	if conditions.RangeStart > 0 {
		db = db.Offset(conditions.RangeStart)
	}

	if conditions.RangeEnd > 0 {
		db = db.Limit(conditions.RangeEnd - conditions.RangeStart)
	}

	var items []AppNode
	results = db.Find(&items)
	return total, items, results.Error

}

func (repo *peerNodeRepository) Save(node AppNode) (AppNode, error) {
	db := sample.DBForProvider(repo.Provider)
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

	db := sample.DBForProvider(repo.Provider)
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

	db := sample.DBForProvider(repo.Provider)
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

	db := sample.DBForProvider(repo.Provider)
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

	db := sample.DBForProvider(repo.Provider)
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

func (repo *peerNodeRepository) TraverseWithPeerIDs(traverse func(AppNode) error, peerIds []string) error {

	db := sample.DBForProvider(repo.Provider)
	if db == nil {
		return sample.ErrSampleDBUnavailable
	}

	rows, err := db.Model(&AppNode{}).Where("peer_id IN ?", peerIds).Rows()
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
