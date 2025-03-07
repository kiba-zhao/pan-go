// Define node search file task repository
package nodesearchfile

import (
	"errors"
	"time"

	appSample "pan/app/sample"
	"pan/app/web"

	"gorm.io/gorm"
)

var ErrNodeSearchTaskNotFound = errors.New("searchtask.NodeSearchTaskRepository Error: Not Found")

type NodeSearchTaskRepository interface {
	// Select retrieves a NodeSearchTask associated with the given ID.
	// It returns the retrieved NodeSearchTask and an error if any issues occur during the query.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
	Select(id uint64) (NodeSearchTask, error)
	// SelectWithQueryAndHash retrieves a NodeSearchTask associated with the given query and hash.
	// It returns the retrieved NodeSearchTask and an error if any issues occur during the query.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
	SelectWithQueryAndHash(query string, hash string) (NodeSearchTask, error)
	// SearchWithStatus retrieves NodeSearchTasks associated with the given status and
	// applying the given range condition for pagination. It returns the total number of
	// results, a slice of NodeSearchTasks, and an error if any issues occur during the
	// query. If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	SearchWithStatus(status uint8, condition web.RangeCondition) (int64, []NodeSearchTask, error)
	// UpdateWithStatus updates the status of the specified NodeSearchTask in the database.
	// It returns the updated NodeSearchTask and an error if any issues occur during the update.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
	UpdateWithStatus(uint8, NodeSearchTask) (NodeSearchTask, error)
	// SelectOrCreate selects a NodeSearchTask associated with the given query and hash,
	// or creates a new one if none exists. It returns the selected or created NodeSearchTask,
	// a boolean indicating whether the NodeSearchTask was created, and an error if any
	// issues occur during the query. If the database is unavailable, it returns
	// appSample.ErrSampleDBUnavailable.
	SelectOrCreate(task NodeSearchTask) (NodeSearchTask, bool, error)
	// Save saves the given NodeSearchTask to the database.
	// It returns the saved NodeSearchTask and an error if any issues occur during the query.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	Save(task NodeSearchTask) (NodeSearchTask, error)
	// SearchWithLifecycle retrieves NodeSearchTasks that have not been updated for the given duration.
	// It returns a slice of NodeSearchTasks and an error if any issues occur during the query.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	SearchWithLifecycle(lifecycle uint64) ([]NodeSearchTask, error)
	// DeleteWithIDs deletes NodeSearchTasks associated with the given IDs.
	// It returns an error if any issues occur during the deletion process.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	DeleteWithIDs(ids ...uint64) error
	// Delete removes the specified NodeSearchTask from the database.
	// It returns an error if any issues occur during the deletion process.
	// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
	// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
	Delete(NodeSearchTask) error
}

type searchTaskRepositoryImpl struct {
	db *gorm.DB
}

// NewNodeSearchTaskRepository creates a new instance of NodeSearchTaskRepository
// using the given Gorm DB instance. If the database connection is nil,
// subsequent repository operations will return ErrNodeSearchTaskNotFound.
func NewNodeSearchTaskRepository(db *gorm.DB) NodeSearchTaskRepository {
	return &searchTaskRepositoryImpl{db: db}
}

// Select retrieves a NodeSearchTask associated with the given ID.
// It returns the retrieved NodeSearchTask and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.

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

// SelectWithQueryAndHash retrieves a NodeSearchTask associated with the given query and hash.
// It returns the retrieved NodeSearchTask and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
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

// SearchWithStatus retrieves NodeSearchTasks associated with the given status and
// applying the given range condition for pagination. It returns the total number of
// results, a slice of NodeSearchTasks, and an error if any issues occur during the
// query. If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
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

// UpdateWithStatus updates the status of the specified NodeSearchTask in the database.
// It returns the updated NodeSearchTask and an error if any issues occur during the update.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
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

// SelectOrCreate selects a NodeSearchTask associated with the given query and status,
// or creates a new one if none exists. It returns the selected or created NodeSearchTask,
// a boolean indicating whether the NodeSearchTask was created, and an error if any
// issues occur during the query. If the database is unavailable, it returns
// appSample.ErrSampleDBUnavailable.
func (repo *searchTaskRepositoryImpl) SelectOrCreate(task NodeSearchTask) (NodeSearchTask, bool, error) {
	db := repo.db
	if db == nil {
		return task, false, appSample.ErrSampleDBUnavailable
	}
	db = db.Where("query = ?", task.Query)
	db = db.Where("status = ?", task.Status)
	db = db.Order("created_at desc")
	results := db.FirstOrCreate(&task)
	return task, results.RowsAffected == 1, results.Error

}

// Save saves the given NodeSearchTask to the database.
// It returns the saved NodeSearchTask and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.
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

// SearchWithLifecycle retrieves NodeSearchTasks that have not been updated for the given lifecycle duration.
// It returns a slice of NodeSearchTasks and an error if any issues occur during the query.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.

func (repo *searchTaskRepositoryImpl) SearchWithLifecycle(lifecycle uint64) ([]NodeSearchTask, error) {
	db := repo.db
	if db == nil {
		return nil, appSample.ErrSampleDBUnavailable
	}
	lifecycleAt := time.Now().Add(-time.Duration(lifecycle))
	var tasks []NodeSearchTask
	results := db.Where("status <> ?", NodeSearchTaskStatusPending).Where("updated_at <= ?", lifecycleAt).Find(&tasks)
	return tasks, results.Error
}

// DeleteWithIDs deletes NodeSearchTasks associated with the given IDs from the database.
// It returns an error if any issues occur during the deletion process.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.

func (repo *searchTaskRepositoryImpl) DeleteWithIDs(ids ...uint64) error {
	db := repo.db
	if db == nil {
		return appSample.ErrSampleDBUnavailable
	}
	results := db.Delete(&NodeSearchTask{}, ids)
	return results.Error
}

// Delete removes the specified NodeSearchTask from the database.
// It returns an error if any issues occur during the deletion process.
// If the database is unavailable, it returns appSample.ErrSampleDBUnavailable.
// If the NodeSearchTask was not found, it returns ErrNodeSearchTaskNotFound.

func (repo *searchTaskRepositoryImpl) Delete(task NodeSearchTask) error {
	db := repo.db
	if db == nil {
		return appSample.ErrSampleDBUnavailable
	}
	results := db.Delete(&task)
	if results.Error == nil && results.RowsAffected != 1 {
		return ErrNodeSearchTaskNotFound
	}
	return results.Error
}
