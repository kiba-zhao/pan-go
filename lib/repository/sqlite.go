package repository

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSqliteDB(dbPath string) (RepositoryDB, error) {

	rootPath := filepath.Dir(dbPath)
	_, err := os.Stat(rootPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(rootPath, 0755)
	}
	if err == nil {
		return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	}
	return nil, err
}
