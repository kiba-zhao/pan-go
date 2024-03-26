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

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var target models.Target
		target.ID = 133
		target.Name = "Target A"
		target.FilePath = "/path_a"
		target.Version = 133
		target.Enabled = true
		target.CreatedAt = time.Now()
		target.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT (.+) FROM `targets`").WithArgs(target.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "file_path", "version", "enabled", "created_at", "updated_at"}).AddRow(
			target.ID, target.Name, target.FilePath, target.Version, target.Enabled, target.CreatedAt, target.UpdatedAt))

		result, err := repo.Select(target.ID, nil)

		assert.Nil(t, err)
		assert.Equal(t, target, result)
	})

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		target := models.Target{
			Name:     "Target A",
			FilePath: "/path_a",
			Enabled:  true,
		}

		mock.ExpectExec("INSERT INTO `targets`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), target.Name, target.FilePath, target.Enabled, target.Invalid, target.Version).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(target, false)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, target.Name, result.Name)
		assert.Equal(t, target.FilePath, result.FilePath)
		assert.Equal(t, target.Enabled, result.Enabled)
		assert.Equal(t, target.Version, result.Version)
		assert.Equal(t, target.Invalid, result.Invalid)

		target.ID = 2
		target.Version = 2
		mock.ExpectExec("UPDATE `targets`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), target.Name, target.FilePath, target.Enabled, target.Invalid, target.Version, target.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(target, false)
		assert.Nil(t, err)
		assert.Equal(t, target.ID, result.ID)
		assert.Equal(t, target.Name, result.Name)
		assert.Equal(t, target.FilePath, result.FilePath)
		assert.Equal(t, target.Enabled, result.Enabled)
		assert.Equal(t, target.Version, result.Version)
		assert.Equal(t, target.Invalid, result.Invalid)

		invalid := true
		target.ID = 3
		target.Version = 3
		target.Invalid = &invalid
		mock.ExpectExec("UPDATE `targets`").WithArgs(sqlmock.AnyArg(), target.Name, target.FilePath, target.Enabled, target.Invalid, target.Version+1, target.Version, target.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(target, true)
		assert.Nil(t, err)
		assert.Equal(t, target.ID, result.ID)
		assert.Equal(t, target.Name, result.Name)
		assert.Equal(t, target.FilePath, result.FilePath)
		assert.Equal(t, target.Enabled, result.Enabled)
		assert.Equal(t, target.Version+1, result.Version)
		assert.Equal(t, target.Invalid, result.Invalid)

	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		target := models.Target{
			ID:       133,
			Version:  13,
			Name:     "Target A",
			FilePath: "/path_a",
			Enabled:  true,
		}

		mock.ExpectExec("UPDATE `targets`").WithArgs(sqlmock.AnyArg(), target.Version, target.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(target)
		assert.Nil(t, err)
	})
}
