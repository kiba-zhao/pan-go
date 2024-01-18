package repositories_test

import (
	"database/sql"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"

	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestPeerRepository ...
func TestPeerRepository(t *testing.T) {

	setup := func() (repo repositories.PeerRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewPeerRepository(db)
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
		hash := []byte{1, 2, 3, 4}
		remoteHash := []byte{11, 12, 13, 14}
		remoteTime := time.Now().Unix()

		mock.ExpectQuery("SELECT (.+) FROM `peers` WHERE (.+) LIMIT 1").WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "enabled", "hash", "remote_hash", "remote_time"}).AddRow(id, enabled, hash, remoteHash, remoteTime))

		row, err := repo.FindOne(id)

		assert.Nil(t, err)
		assert.Equal(t, id, row.ID)
		assert.Equal(t, enabled, row.Enabled)
		assert.Equal(t, hash, row.Hash)
		assert.Equal(t, remoteHash, row.RemoteHash)
		assert.Equal(t, remoteTime, row.RemoteTime)

	})

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var peer models.Peer
		peer.ID = uuid.New().String()
		peer.Enabled = true
		peer.Hash = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
		peer.RemoteHash = []byte{1, 2, 3, 4, 5}
		peer.RemoteTime = time.Now().Unix()

		mock.ExpectExec("UPDATE `peers` SET (.+) WHERE (.+) `id` =").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), peer.Enabled, peer.Hash, peer.RemoteHash, peer.RemoteTime, peer.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.Save(peer)

		assert.Nil(t, err)

	})
}
