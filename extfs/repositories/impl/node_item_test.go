package impl_test

import (
	"database/sql"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"pan/extfs/repositories/impl"
	mockedApp "pan/mocks/pan/app"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNodeItemRepo(t *testing.T) {
	setup := func() (repositories.NodeItemRepository, *sql.DB, sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		provider := new(mockedApp.MockRepositoryDBProvider)
		provider.On("DB").Return(db)
		repo := &impl.NodeItemRepository{Provider: provider}
		return repo, mockDB, mock
	}

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item models.NodeItem
		item.ID = 123
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `node_items`").WithArgs(item.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "file_path", "enabled", "created_at", "updated_at"}).AddRow(item.ID, item.Name, item.FilePath, item.Enabled, item.CreatedAt, item.UpdatedAt))

		result, err := repo.Select(item.ID)
		assert.Nil(t, err)
		assert.Equal(t, item, result)
	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item models.NodeItem
		item.ID = 123
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		mock.ExpectExec("UPDATE `node_items`").WithArgs(sqlmock.AnyArg(), item.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(item)
		assert.Nil(t, err)
	})

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item models.NodeItem
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		mock.ExpectExec("INSERT INTO `node_items`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), item.Name, item.FilePath, item.Enabled).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(item)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, item.Name, result.Name)
		assert.Equal(t, item.FilePath, result.FilePath)
		assert.Equal(t, *item.Enabled, *result.Enabled)

		item.ID = 123
		mock.ExpectExec("UPDATE `node_items`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), item.Name, item.FilePath, item.Enabled, item.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(item)
		assert.Nil(t, err)
		assert.Equal(t, item.ID, result.ID)
		assert.Equal(t, item.Name, result.Name)
		assert.Equal(t, item.FilePath, result.FilePath)
		assert.Equal(t, *item.Enabled, *result.Enabled)
	})
}
