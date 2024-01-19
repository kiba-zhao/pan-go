package repositories_test

import (
	"database/sql"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestExtFS(t *testing.T) {

	// setup function to repositories.ExtFSRepository
	setup := func() (repo repositories.ExtFSRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewExtFSRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("GetLatestOne", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var extFSRow models.ExtFS
		extFSRow.Hash = []byte("hash")
		extFSRow.CreatedAt = time.Now()
		mock.ExpectQuery("SELECT  (.+) FROM `ext_fs` (.+) LIMIT 1").WillReturnRows(sqlmock.NewRows([]string{"hash", "created_at"}).AddRow(extFSRow.Hash, extFSRow.CreatedAt))

		row, err := repo.GetLatestOne()

		assert.Nil(t, err)
		assert.Equal(t, extFSRow.Hash, row.Hash)
		assert.Equal(t, extFSRow.CreatedAt, row.CreatedAt)
	})
}
