package appnode_test

import (
	"database/sql"
	"encoding/base64"
	appnode "pan/app/app_node"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAppNodeRepo(t *testing.T) {

	setup := func() (repo appnode.AppNodeRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = appnode.NewAppNodeRepository(db)
		return
	}

	t.Run("SelectByPeerID", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.ID = 123
		peerNode.Name = "peer node name"
		peerNode.PeerID = base64.StdEncoding.EncodeToString([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `app_nodes`").WithArgs(peerNode.PeerID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "peer_id", "blocked", "created_at", "updated_at"}).AddRow(peerNode.ID, peerNode.Name, peerNode.PeerID, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt))

		result, err := repo.SelectByPeerID(peerNode.PeerID)
		assert.Nil(t, err)
		assert.Equal(t, peerNode, result)
	})

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.ID = 123
		peerNode.Name = "peer node name"
		peerNode.PeerID = base64.StdEncoding.EncodeToString([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `app_nodes`").WithArgs(peerNode.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "peer_id", "blocked", "created_at", "updated_at"}).AddRow(peerNode.ID, peerNode.Name, peerNode.PeerID, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt))

		result, err := repo.Select(peerNode.ID)
		assert.Nil(t, err)
		assert.Equal(t, peerNode, result)

	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.ID = 123
		peerNode.Name = "peer node name"
		peerNode.PeerID = base64.StdEncoding.EncodeToString([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectExec("UPDATE `app_nodes`").WithArgs(sqlmock.AnyArg(), peerNode.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(peerNode)
		assert.Nil(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.ID = 123
		peerNode.Name = "peer node name"
		peerNode.PeerID = base64.StdEncoding.EncodeToString([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()
		rowCount := int64(123)

		mock.ExpectQuery("SELECT .* FROM `app_nodes`").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow(rowCount))
		mock.ExpectQuery("SELECT .* FROM `app_nodes`").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "peer_id", "blocked", "created_at", "updated_at"}).AddRow(peerNode.ID, peerNode.Name, peerNode.PeerID, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt))

		total, rows, err := repo.Search(appnode.AppNodeSearchCondition{})
		assert.Nil(t, err)
		assert.Equal(t, rowCount, total)
		assert.Equal(t, 1, len(rows))
		assert.Equal(t, peerNode, rows[0])
	})

	t.Run("Save", func(t *testing.T) {

		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.Name = "peer node name"
		peerNode.PeerID = base64.StdEncoding.EncodeToString([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectExec("INSERT INTO `app_nodes`").WithArgs(peerNode.PeerID, peerNode.Name, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(peerNode)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, peerNode.Name, result.Name)
		assert.Equal(t, peerNode.PeerID, result.PeerID)
		assert.Equal(t, peerNode.Blocked, result.Blocked)

		peerNode.ID = 123

		mock.ExpectExec("UPDATE `app_nodes`").WithArgs(peerNode.PeerID, peerNode.Name, peerNode.Blocked, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(peerNode)
		assert.Nil(t, err)
		assert.Equal(t, peerNode.ID, result.ID)
		assert.Equal(t, peerNode.Name, result.Name)
		assert.Equal(t, peerNode.PeerID, result.PeerID)
		assert.Equal(t, peerNode.Blocked, result.Blocked)
	})
}
