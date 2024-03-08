package repositories_test

import (
	"database/sql"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTargetRepository(t *testing.T) {

	setup := func() (repo repositories.TargetRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewTargetRepository(db)
		return

	}

	t.Run("Search", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		rowCount := int64(123)
		var target models.Target
		target.ID = 1
		target.Name = "Target A"
		target.FilePath = "/path_a"
		target.Version = 133
		target.Enabled = true
		target.CreatedAt = time.Now()
		target.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT (.+) FROM `targets`").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(rowCount))
		mock.ExpectQuery("SELECT (.+) FROM `targets`").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "file_path", "version", "enabled", "created_at", "updated_at"}).AddRow(
			target.ID, target.Name, target.FilePath, target.Version, target.Enabled, target.CreatedAt, target.UpdatedAt))

		total, rows, err := repo.Search(models.TargetSearchCondition{})

		assert.Nil(t, err)
		assert.Equal(t, rowCount, total)
		assert.Equal(t, 1, len(rows))
		assert.Equal(t, target, rows[0])
	})

}
