package nodeitem_test

import (
	"database/sql"
	nodeitem "pan/extfs/node_item"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNodeItemRepo(t *testing.T) {
	setup := func() (nodeitem.NodeItemRepository, *sql.DB, sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo := nodeitem.NewNodeItemRepository(db)
		return repo, mockDB, mock
	}

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item nodeitem.NodeItem
		item.ID = 123
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.FileType = "node item file type"
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `node_items`").WithArgs(item.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "file_path", "file_type", "enabled", "created_at", "updated_at"}).AddRow(item.ID, item.Name, item.FilePath, item.FileType, item.Enabled, item.CreatedAt, item.UpdatedAt))

		result, err := repo.Select(item.ID)
		assert.Nil(t, err)
		assert.Equal(t, item, result)
	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item nodeitem.NodeItem
		item.ID = 123
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.FileType = "node item file type"
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
		var item nodeitem.NodeItem
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.FileType = nodeitem.FileTypeFolder
		item.MimeType = "node item mime type"
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		mock.ExpectExec("INSERT INTO `node_items`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), item.Name, item.FilePath, item.FileType, item.MimeType, item.Enabled).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(item)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, item.Name, result.Name)
		assert.Equal(t, item.FilePath, result.FilePath)
		assert.Equal(t, item.FileType, result.FileType)
		assert.Equal(t, item.MimeType, result.MimeType)
		assert.Equal(t, *item.Enabled, *result.Enabled)

		item.ID = 123
		mock.ExpectExec("UPDATE `node_items`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), item.Name, item.FilePath, item.FileType, item.MimeType, item.Enabled, item.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(item)
		assert.Nil(t, err)
		assert.Equal(t, item.ID, result.ID)
		assert.Equal(t, item.Name, result.Name)
		assert.Equal(t, item.FilePath, result.FilePath)
		assert.Equal(t, item.FileType, result.FileType)
		assert.Equal(t, item.MimeType, result.MimeType)
		assert.Equal(t, *item.Enabled, *result.Enabled)
	})

	t.Run("TraverseAll", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		enabled := true
		var item nodeitem.NodeItem
		item.Name = "node item name"
		item.FilePath = "node item file path"
		item.FileType = nodeitem.FileTypeFolder
		item.Enabled = &enabled
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()

		var nextItem nodeitem.NodeItem
		nextItem.Name = "next node item name"
		nextItem.FilePath = "next node item file path"
		nextItem.FileType = nodeitem.FileTypeFolder
		nextItem.Enabled = &enabled
		nextItem.CreatedAt = time.Now()
		nextItem.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `node_items`").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "file_path", "file_type", "enabled", "created_at", "updated_at"}).AddRow(item.ID, item.Name, item.FilePath, item.FileType, item.Enabled, item.CreatedAt, item.UpdatedAt).AddRow(nextItem.ID, nextItem.Name, nextItem.FilePath, nextItem.FileType, nextItem.Enabled, nextItem.CreatedAt, nextItem.UpdatedAt))

		times := 0
		err := repo.TraverseAll(func(model nodeitem.NodeItem) error {
			if times == 0 {
				assert.Equal(t, item, model)
			}
			if times == 1 {
				assert.Equal(t, nextItem, model)
			}
			times++
			return nil
		})
		assert.Nil(t, err)
		assert.Equal(t, 2, times)
	})
}
