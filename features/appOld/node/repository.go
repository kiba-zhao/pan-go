package node

import (
	"errors"
	"pan/lib/feature"
	"pan/lib/repository"
	"pan/lib/web"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAppNodeNotFound = errors.New("appnode.AppNodeRepository Error: Not Found")
var ErrAppNodeInvaild = errors.New("appnode.AppNodeRepository Error: Invaild")

type AppNodeRepository interface {
	repository.Repository
	Search(AppNodeSearchCondition) (int64, []AppNode, error)
	Create(AppNode) (AppNode, error)
	Update(AppNode) (AppNode, error)
	Select(uint) (AppNode, error)
	SelectByName(string) (AppNode, error)
	Delete(AppNode) error
	SelectByPeerID(string) (AppNode, error)
	TraverseAll(func(AppNode) error) error
}

type peerNodeRepository struct {
	feature.Repository
}

func NewAppNodeRepository() AppNodeRepository {
	return &peerNodeRepository{}
}

var _ = (AppNodeRepository)((*peerNodeRepository)(nil))

func (repo *peerNodeRepository) Search(conditions AppNodeSearchCondition) (int64, []AppNode, error) {
	db := repo.DB()
	if db == nil {
		return 0, nil, feature.ErrFeatureRepositoryDBUnavailable
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

func (repo *peerNodeRepository) Create(node AppNode) (AppNode, error) {
	db := repo.DB()
	if db == nil {
		return node, feature.ErrFeatureRepositoryDBUnavailable
	}
	if node.ID > 0 {
		return node, ErrAppNodeInvaild
	}
	results := db.Create(&node)
	err := results.Error
	if err == nil && results.RowsAffected < 1 {
		err = ErrAppNodeNotFound
	}
	return node, err
}

func (repo *peerNodeRepository) Update(node AppNode) (AppNode, error) {
	db := repo.DB()
	if db == nil {
		return node, feature.ErrFeatureRepositoryDBUnavailable
	}
	if node.ID <= 0 {
		return node, ErrAppNodeInvaild
	}

	db = db.Begin()

	// remove NetworkAddr
	tx := db.Where("app_node_id = ?", node.ID)
	if len(node.NetworkAddrs) > 0 {
		addrIds := make([]uint64, 0)
		for _, addr := range node.NetworkAddrs {
			if addr.ID > 0 {
				addrIds = append(addrIds, addr.ID)
			}
		}
		if len(addrIds) > 0 {
			tx = tx.Not(addrIds)
		}
	}
	tx = tx.Delete(&NetworkAddr{})
	err := tx.Error

	if err == nil {
		results := db.Save(&node)
		err = results.Error
		if err == nil && results.RowsAffected < 1 {
			err = ErrAppNodeNotFound
		}
	}

	if err == nil {
		results := db.Commit()
		err = results.Error
	}
	if err != nil {
		db.Rollback()
	}

	return node, err
}

func (repo *peerNodeRepository) Select(id uint) (AppNode, error) {

	db := repo.DB()
	if db == nil {
		return AppNode{}, feature.ErrFeatureRepositoryDBUnavailable
	}
	var peerNode AppNode

	results := db.Model(&peerNode).Preload("NetworkAddrs").Take(&peerNode, id)
	if results.Error == gorm.ErrRecordNotFound {
		return peerNode, ErrAppNodeNotFound
	}
	return peerNode, results.Error
}

func (repo *peerNodeRepository) SelectByName(name string) (AppNode, error) {

	db := repo.DB()
	if db == nil {
		return AppNode{}, feature.ErrFeatureRepositoryDBUnavailable
	}
	var peerNode AppNode
	results := db.Where("name = ?", name).Take(&peerNode)
	if results.Error == gorm.ErrRecordNotFound {
		return peerNode, ErrAppNodeNotFound
	}
	return peerNode, results.Error
}

func (repo *peerNodeRepository) Delete(node AppNode) error {

	db := repo.DB()
	if db == nil {
		return feature.ErrFeatureRepositoryDBUnavailable
	}
	results := db.Select(clause.Associations).Delete(&node)
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrAppNodeNotFound
	}
	return results.Error
}

func (repo *peerNodeRepository) SelectByPeerID(peerId string) (AppNode, error) {

	db := repo.DB()
	if db == nil {
		return AppNode{}, feature.ErrFeatureRepositoryDBUnavailable
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

	db := repo.DB()
	if db == nil {
		return feature.ErrFeatureRepositoryDBUnavailable
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
