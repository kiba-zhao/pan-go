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

func TestFilesStateRepository(t *testing.T) {

	setup := func() (repo repositories.FilesStateRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewFilesStateRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("GetLastOne", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var filesStateRow models.FilesState
		filesStateRow.Hash = []byte("hash")
		filesStateRow.Origin = 1
		filesStateRow.CreatedAt = time.Now()

		mock.ExpectQuery("SELECT  (.+) FROM `files_states` (.+) LIMIT 1").WillReturnRows(sqlmock.NewRows([]string{"hash", "origin", "created_at"}).AddRow(filesStateRow.Hash, filesStateRow.Origin, filesStateRow.CreatedAt))

		row, err := repo.GetLastOne()

		assert.Nil(t, err)
		assert.Equal(t, filesStateRow.Hash, row.Hash)
		assert.Equal(t, filesStateRow.Origin, row.Origin)
		assert.Equal(t, filesStateRow.CreatedAt, row.CreatedAt)
	})

}
