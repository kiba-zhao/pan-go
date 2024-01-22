package repositories_test

import (
	"database/sql"
	"pan/modules/extfs/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestRemotePeerRepository ...
func TestRemotePeerRepository(t *testing.T) {

	// setup function to repositories.RemotePeerRepository
	setup := func() (repo repositories.RemotePeerRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewRemotePeerRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("FindOne", func(t *testing.T) {

		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		id := uuid.New().String()
		enabled := true

		mock.ExpectQuery("SELECT (.+) FROM `remote_peers` WHERE (.+) LIMIT 1").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "enabled"}).AddRow(id, enabled))

		row, err := repo.FindOne(id)

		assert.Nil(t, err)
		assert.Equal(t, id, row.ID)
		assert.Equal(t, enabled, row.Enabled)

	})
}
