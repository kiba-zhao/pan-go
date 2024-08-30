package impl

import (
	"pan/app/constant"
	"pan/app/models"
	"pan/app/repositories"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NodeRepository struct {
	Provider repositories.RepositoryDBProvider
}

func (repo *NodeRepository) Search(conditions models.NodeSearchCondition) (int64, []models.Node, error) {

	db := repositories.DBForProvider(repo.Provider)
	if db == nil {
		return 0, nil, constant.ErrUnavailable
	}

	if len(conditions.Keyword) > 0 {
		tx := db
		keywords := strings.Split(conditions.Keyword, ",")
		for _, keyword := range keywords {
			trimKeyword := strings.Trim(keyword, " ")
			if len(trimKeyword) > 0 {
				tx = tx.Or("name like ?", "%"+keyword+"%")
				tx = tx.Or("nodeId = ?", keyword)
			}
		}
		db = db.Where(tx)
	}

	if conditions.Blocked != nil {
		db = db.Where("blocked = ?", *conditions.Blocked)
	}

	total := int64(0)
	results := db.Model(&models.Node{}).Count(&total)

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

	var items []models.Node
	results = db.Find(&items)
	return total, items, results.Error

}

func (repo *NodeRepository) Save(node models.Node) (models.Node, error) {
	db := repositories.DBForProvider(repo.Provider)
	if db == nil {
		return node, constant.ErrUnavailable
	}

	results := db.Save(&node)
	if results.Error == nil && results.RowsAffected != 1 {
		return node, constant.ErrNotFound
	}
	return node, results.Error
}

func (repo *NodeRepository) Select(id uint) (models.Node, error) {

	db := repositories.DBForProvider(repo.Provider)
	if db == nil {
		return models.Node{}, constant.ErrUnavailable
	}
	var node models.Node
	results := db.Take(&node, id)
	if results.Error == gorm.ErrRecordNotFound {
		return node, constant.ErrNotFound
	}
	return node, results.Error
}

func (repo *NodeRepository) Delete(node models.Node) error {

	db := repositories.DBForProvider(repo.Provider)
	if db == nil {
		return constant.ErrUnavailable
	}
	results := db.Delete(&node)
	if results.Error == nil && results.RowsAffected != 1 {
		return constant.ErrNotFound
	}
	return results.Error
}

func (repo *NodeRepository) SelectByNodeID(nodeId string) (models.Node, error) {

	db := repositories.DBForProvider(repo.Provider)
	if db == nil {
		return models.Node{}, constant.ErrUnavailable
	}
	var node models.Node
	node.NodeID = nodeId
	results := db.Where(&node).Take(&node)
	if results.Error == gorm.ErrRecordNotFound {
		return node, constant.ErrNotFound
	}
	return node, results.Error
}
