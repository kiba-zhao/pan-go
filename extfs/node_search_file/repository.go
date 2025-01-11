package nodesearchfile

import (
	"errors"
	"pan/app/web"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var ErrNodeSearchFileDBUnavailable = errors.New("searchfile.NodeSearchFileRepository Error: Database unavailable")
var ErrNodeSearchFileNotFound = errors.New("searchfile.NodeSearchFileRepository Error: Not Found")

func (NodeSearchFile) TableName() string {
	return "node_file_rates"
}

func generateTableNameWithNodeSearchFile(searchFile *NodeSearchFile, taskId uint64) string {
	return strings.Join([]string{searchFile.TableName(), strconv.FormatUint(taskId, 10)}, "_")
}

func NodeSearchFileWithTaskID(searchFile *NodeSearchFile, taskId uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tableName := generateTableNameWithNodeSearchFile(searchFile, taskId)
		return db.Table(tableName)
	}
}

type NodeSearchFileRepository interface {
	Init(taskId uint64) error
	Destroy(taskId uint64) error
	Search(taskId uint64, condition web.RangeCondition) (int64, []NodeSearchFile, error)
	Save(taskId uint64, searchFile NodeSearchFile) (NodeSearchFile, error)
}

type searchFileRepositoryImpl struct {
	db *gorm.DB
}

func NewNodeSearchFileRepository(db *gorm.DB) NodeSearchFileRepository {
	return &searchFileRepositoryImpl{db: db}
}

func (repo *searchFileRepositoryImpl) Init(taskId uint64) error {
	db := repo.db
	if db == nil {
		return ErrNodeSearchFileDBUnavailable
	}
	var model NodeSearchFile
	db = db.Scopes(NodeSearchFileWithTaskID(&model, taskId))
	migrator := db.Migrator()
	if migrator.HasTable(&model) {
		return nil
	}
	return migrator.CreateTable(&model)
}

func (repo *searchFileRepositoryImpl) Destroy(taskId uint64) error {
	db := repo.db
	if db == nil {
		return ErrNodeSearchFileDBUnavailable
	}
	var model NodeSearchFile
	db = db.Scopes(NodeSearchFileWithTaskID(&model, taskId))
	migrator := db.Migrator()
	if migrator.HasTable(&model) {
		return migrator.DropTable(&model)
	}
	return nil
}

func (repo *searchFileRepositoryImpl) Search(taskId uint64, condition web.RangeCondition) (int64, []NodeSearchFile, error) {
	db := repo.db
	if db == nil {
		return 0, nil, ErrNodeSearchFileDBUnavailable
	}
	var model NodeSearchFile
	db = db.Scopes(NodeSearchFileWithTaskID(&model, taskId))

	total := int64(0)
	results := db.Model(&NodeSearchFile{}).Count(&total)

	if results.Error != nil || total <= 0 {
		return total, nil, results.Error
	}

	db = db.Scopes(web.PaginateWithRangeCondition(&condition))
	var searchFiles []NodeSearchFile

	results = db.Find(&searchFiles)
	return total, searchFiles, results.Error
}

func (repo *searchFileRepositoryImpl) Save(taskId uint64, searchFile NodeSearchFile) (NodeSearchFile, error) {
	db := repo.db
	if db == nil {
		return searchFile, ErrNodeSearchFileDBUnavailable
	}

	db = db.Scopes(NodeSearchFileWithTaskID(&searchFile, taskId))
	results := db.Save(&searchFile)
	if results.Error == nil && results.RowsAffected != 1 {
		return searchFile, ErrNodeSearchFileNotFound
	}
	return searchFile, results.Error
}
