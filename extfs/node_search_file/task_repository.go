package nodesearchfile

import (
	"errors"

	appSample "pan/app/sample"
	"pan/app/web"

	"gorm.io/gorm"
)

var ErrNodeSearchTaskNotFound = errors.New("searchtask.NodeSearchTaskRepository Error: Not Found")

type NodeSearchTaskRepository interface {
	Select(id uint64) (NodeSearchTask, error)
	SelectWithQueryAndHash(query string, hash string) (NodeSearchTask, error)
	SearchWithStatus(status uint8, condition web.RangeCondition) (int64, []NodeSearchTask, error)
	UpdateWithStatus(uint8, NodeSearchTask) (NodeSearchTask, error)
	SelectOrCreate(task NodeSearchTask) (NodeSearchTask, bool, error)
	Save(task NodeSearchTask) (NodeSearchTask, error)
}

type searchTaskRepositoryImpl struct {
	db *gorm.DB
}

func NewNodeSearchTaskRepository(db *gorm.DB) NodeSearchTaskRepository {
	return &searchTaskRepositoryImpl{db: db}
}

func (repo *searchTaskRepositoryImpl) Select(id uint64) (NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return NodeSearchTask{}, appSample.ErrSampleDBUnavailable
	}
	var task NodeSearchTask
	results := db.Take(&task, id)
	if results.Error == gorm.ErrRecordNotFound {
		return task, ErrNodeSearchTaskNotFound
	}
	return task, results.Error
}

func (repo *searchTaskRepositoryImpl) SelectWithQueryAndHash(query string, hash string) (NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return NodeSearchTask{}, appSample.ErrSampleDBUnavailable
	}
	var task NodeSearchTask
	db = db.Where("query = ?", query)
	if len(hash) > 0 {
		db = db.Where("hash = ?", hash)
	}
	results := db.Order("created_at desc").Take(&task)
	if results.Error == gorm.ErrRecordNotFound {
		return task, ErrNodeSearchTaskNotFound
	}
	return task, results.Error
}

func (repo *searchTaskRepositoryImpl) SearchWithStatus(status uint8, condition web.RangeCondition) (int64, []NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return 0, nil, appSample.ErrSampleDBUnavailable
	}

	total := int64(0)
	results := db.Model(&NodeSearchTask{}).Count(&total)

	if results.Error != nil || total <= 0 {
		return total, nil, results.Error
	}

	db = db.Scopes(web.PaginateWithRangeCondition(&condition))
	var searchTasks []NodeSearchTask

	results = db.Order("created_at desc").Where("status = ?", status).Find(&searchTasks)
	return total, searchTasks, results.Error
}

func (repo *searchTaskRepositoryImpl) UpdateWithStatus(status uint8, task NodeSearchTask) (NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return task, appSample.ErrSampleDBUnavailable
	}
	results := db.Model(&task).Where("id = ?", task.ID).Where("status = ?", status).Limit(1).Updates(&task)
	if results.Error == nil && results.RowsAffected != 1 {
		return task, ErrNodeSearchTaskNotFound
	}
	return task, results.Error
}

func (repo *searchTaskRepositoryImpl) SelectOrCreate(task NodeSearchTask) (NodeSearchTask, bool, error) {
	db := repo.db
	if db == nil {
		return task, false, appSample.ErrSampleDBUnavailable
	}
	db = db.Where("query = ?", task.Query)
	db = db.Where("status = ?", task.Status)
	db = db.Order("created_at desc")
	results := db.FirstOrCreate(&task)
	if results.Error == nil {
		return task, results.RowsAffected == 1, results.Error
	}
	return task, false, results.Error
}

func (repo *searchTaskRepositoryImpl) Save(task NodeSearchTask) (NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return task, appSample.ErrSampleDBUnavailable
	}
	results := db.Save(&task)
	if results.Error == nil && results.RowsAffected != 1 {
		return task, ErrNodeSearchTaskNotFound
	}
	return task, results.Error
}
