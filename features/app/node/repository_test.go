package node_test

import (
	"database/sql"
	appnode "pan/features/app/node"
	"pan/lib/peer"
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
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
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
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.NetworkAddrs = []appnode.NetworkAddr{
			{ID: 1, AppNodeID: peerNode.ID, Address: "address1"},
		}
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT .* FROM `app_nodes`").WithArgs(peerNode.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "peer_id", "blocked", "created_at", "updated_at"}).AddRow(peerNode.ID, peerNode.Name, peerNode.PeerID, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt))
		mock.ExpectQuery("SELECT .* FROM `network_addrs`").WithArgs(peerNode.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "app_node_id", "address"}).AddRow(peerNode.NetworkAddrs[0].ID, peerNode.NetworkAddrs[0].AppNodeID, peerNode.NetworkAddrs[0].Address))

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
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		mock.ExpectExec("DELETE FROM `network_addrs`").WithArgs(peerNode.ID).WillReturnResult(sqlmock.NewResult(1, 1))
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
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
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

	t.Run("Create", func(t *testing.T) {

		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.Name = "peer node name"
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.NetworkAddrs = []appnode.NetworkAddr{
			{Address: "address1"},
		}
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		peerNodeId := int64(333)
		networkAddrId := int64(444)
		mock.ExpectExec("INSERT INTO `app_nodes`").WithArgs(peerNode.PeerID, peerNode.Name, peerNode.Blocked, peerNode.CreatedAt, peerNode.UpdatedAt, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(peerNodeId, 1))
		mock.ExpectExec("INSERT INTO `network_addrs`").WithArgs(peerNodeId, peerNode.NetworkAddrs[0].Address, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(networkAddrId, 1))

		result, err := repo.Create(peerNode)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint(0))
		assert.Equal(t, peerNode.Name, result.Name)
		assert.Equal(t, peerNode.PeerID, result.PeerID)
		assert.Equal(t, peerNode.Blocked, result.Blocked)
		assert.Equal(t, uint64(networkAddrId), result.NetworkAddrs[0].ID)
		assert.Equal(t, uint(peerNodeId), result.NetworkAddrs[0].AppNodeID)
		assert.Equal(t, peerNode.NetworkAddrs[0].Address, result.NetworkAddrs[0].Address)

		// peerNode.ID = 123
		// peerNode.NetworkAddrs = []appnode.NetworkAddr{
		// 	{Address: "address1"},
		// 	{ID: 345, AppNodeID: peerNode.ID, Address: "address2"},
		// }

		// networkAddrId = int64(555)
		// mock.ExpectExec("UPDATE `app_nodes`").WithArgs(peerNode.PeerID, peerNode.Name, peerNode.Blocked, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		// mock.ExpectExec("INSERT INTO `network_addrs`").WithArgs(
		// 	peerNode.ID, peerNode.NetworkAddrs[0].Address, sqlmock.AnyArg(), sqlmock.AnyArg(),
		// 	peerNode.ID, peerNode.NetworkAddrs[1].Address, sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.NetworkAddrs[1].ID,
		// ).WillReturnResult(sqlmock.NewResult(networkAddrId, 1))

		// result, err = repo.Update(peerNode)
		// assert.Nil(t, err)
		// assert.Equal(t, peerNode.ID, result.ID)
		// assert.Equal(t, peerNode.Name, result.Name)
		// assert.Equal(t, peerNode.PeerID, result.PeerID)
		// assert.Equal(t, peerNode.Blocked, result.Blocked)
		// assert.Equal(t, uint64(networkAddrId), result.NetworkAddrs[0].ID)
		// assert.Equal(t, peerNode.ID, result.NetworkAddrs[0].AppNodeID)
		// assert.Equal(t, peerNode.NetworkAddrs[0].Address, result.NetworkAddrs[0].Address)
		// assert.Equal(t, uint64(peerNode.NetworkAddrs[1].ID), result.NetworkAddrs[1].ID)
		// assert.Equal(t, peerNode.ID, result.NetworkAddrs[1].AppNodeID)
		// assert.Equal(t, peerNode.NetworkAddrs[1].Address, result.NetworkAddrs[1].Address)
	})

	t.Run("Update", func(t *testing.T) {

		repo, mockDB, mock := setup()
		defer mockDB.Close()

		var peerNode appnode.AppNode
		peerNode.ID = 123
		peerNode.Name = "peer node name"
		peerNode.PeerID = peer.EncodePeerID([]byte("peer node id"))
		peerNode.Blocked = true
		peerNode.NetworkAddrs = []appnode.NetworkAddr{
			{Address: "address1"},
			{ID: 345, AppNodeID: peerNode.ID, Address: "address2"},
			{ID: 346, AppNodeID: peerNode.ID, Address: "address3"},
		}
		peerNode.CreatedAt = time.Now()
		peerNode.UpdatedAt = time.Now()

		networkAddrId := int64(555)
		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM `network_addrs`").WithArgs(peerNode.ID, peerNode.NetworkAddrs[1].ID, peerNode.NetworkAddrs[2].ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("UPDATE `app_nodes`").WithArgs(peerNode.PeerID, peerNode.Name, peerNode.Blocked, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.ID).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("INSERT INTO `network_addrs`").WithArgs(
			peerNode.ID, peerNode.NetworkAddrs[0].Address, sqlmock.AnyArg(), sqlmock.AnyArg(),
			peerNode.ID, peerNode.NetworkAddrs[1].Address, sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.NetworkAddrs[1].ID,
			peerNode.ID, peerNode.NetworkAddrs[2].Address, sqlmock.AnyArg(), sqlmock.AnyArg(), peerNode.NetworkAddrs[2].ID,
		).WillReturnResult(sqlmock.NewResult(networkAddrId, 1))
		mock.ExpectCommit()

		result, err := repo.Update(peerNode)
		assert.Nil(t, err)
		assert.Equal(t, peerNode.ID, result.ID)
		assert.Equal(t, peerNode.Name, result.Name)
		assert.Equal(t, peerNode.PeerID, result.PeerID)
		assert.Equal(t, peerNode.Blocked, result.Blocked)
		assert.Equal(t, uint64(networkAddrId), result.NetworkAddrs[0].ID)
		assert.Equal(t, peerNode.ID, result.NetworkAddrs[0].AppNodeID)
		assert.Equal(t, peerNode.NetworkAddrs[0].Address, result.NetworkAddrs[0].Address)
		assert.Equal(t, uint64(peerNode.NetworkAddrs[1].ID), result.NetworkAddrs[1].ID)
		assert.Equal(t, peerNode.ID, result.NetworkAddrs[1].AppNodeID)
		assert.Equal(t, peerNode.NetworkAddrs[1].Address, result.NetworkAddrs[1].Address)
		assert.Equal(t, uint64(peerNode.NetworkAddrs[2].ID), result.NetworkAddrs[2].ID)
		assert.Equal(t, peerNode.ID, result.NetworkAddrs[2].AppNodeID)
		assert.Equal(t, peerNode.NetworkAddrs[2].Address, result.NetworkAddrs[2].Address)
	})
}
