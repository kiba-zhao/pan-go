package impl_test

import (
	"database/sql"
	"encoding/base64"
	"pan/app/models"
	"pan/app/repositories"
	"pan/app/repositories/impl"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	mockedApp "pan/mocks/pan/app"
)

func TestNodeRepo(t *testing.T) {

	setup := func() (repo repositories.NodeRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
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
		repo = &impl.NodeRepository{Provider: provider}
		return
	}

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var node models.Node
		node.ID = 123
		node.Name = "node name"
		node.NodeID = base64.StdEncoding.EncodeToString([]byte("node id"))
		node.Blocked = true
		node.CreatedAt = time.Now()
		node.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `nodes`").WithArgs(node.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "node_id", "blocked", "created_at", "updated_at"}).AddRow(node.ID, node.Name, node.NodeID, node.Blocked, node.CreatedAt, node.UpdatedAt))

		result, err := repo.Select(node.ID)
		assert.Nil(t, err)
		assert.Equal(t, node, result)

	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var node models.Node
		node.ID = 123
		node.Name = "node name"
		node.NodeID = base64.StdEncoding.EncodeToString([]byte("node id"))
		node.Blocked = true
		node.CreatedAt = time.Now()
		node.UpdatedAt = time.Now()

		mock.ExpectExec("UPDATE `nodes`").WithArgs(sqlmock.AnyArg(), node.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(node)
		assert.Nil(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var node models.Node
		node.ID = 123
		node.Name = "node name"
		node.NodeID = base64.StdEncoding.EncodeToString([]byte("node id"))
		node.Blocked = true
		node.CreatedAt = time.Now()
		node.UpdatedAt = time.Now()
		rowCount := int64(123)

		mock.ExpectQuery("SELECT .* FROM `nodes`").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(rowCount))
		mock.ExpectQuery("SELECT .* FROM `nodes`").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "node_id", "blocked", "created_at", "updated_at"}).AddRow(node.ID, node.Name, node.NodeID, node.Blocked, node.CreatedAt, node.UpdatedAt))

		total, rows, err := repo.Search(models.NodeSearchCondition{})
		assert.Nil(t, err)
		assert.Equal(t, rowCount, total)
		assert.Equal(t, 1, len(rows))
		assert.Equal(t, node, rows[0])
	})

	t.Run("Save", func(t *testing.T) {

		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var node models.Node
		node.Name = "node name"
		node.NodeID = base64.StdEncoding.EncodeToString([]byte("node id"))
		node.Blocked = true
		node.CreatedAt = time.Now()
		node.UpdatedAt = time.Now()

		mock.ExpectExec("INSERT INTO `nodes`").WithArgs(node.NodeID, node.Name, node.Blocked, node.CreatedAt, node.UpdatedAt, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(node)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, node.Name, result.Name)
		assert.Equal(t, node.NodeID, result.NodeID)
		assert.Equal(t, node.Blocked, result.Blocked)

		node.ID = 123

		mock.ExpectExec("UPDATE `nodes`").WithArgs(node.NodeID, node.Name, node.Blocked, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), node.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(node)
		assert.Nil(t, err)
		assert.Equal(t, node.ID, result.ID)
		assert.Equal(t, node.Name, result.Name)
		assert.Equal(t, node.NodeID, result.NodeID)
		assert.Equal(t, node.Blocked, result.Blocked)
	})
}
