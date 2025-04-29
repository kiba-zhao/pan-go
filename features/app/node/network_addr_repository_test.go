package node_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	appnode "pan/features/app/node"
)

func TestNetworkAddrRepository(t *testing.T) {

	setup := func() (repo appnode.NetworkAddrRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = appnode.NewNetworkAddrRepository(db)
		return
	}

	t.Run("SeqByPeerId", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		peerId := "peerId"
		var entity appnode.NetworkAddr
		entity.AppNodeID = 123
		entity.Address = "address"
		entity.ID = 123

		mock.ExpectQuery("SELECT (.+) FROM `network_addrs` INNER JOIN `app_nodes` `AppNode` ON `network_addrs`.`app_node_id` = `AppNode`.`id`").WithArgs(peerId).WillReturnRows(sqlmock.NewRows([]string{"id", "app_node_id", "address", "created_at", "updated_at"}).AddRow(entity.ID, entity.AppNodeID, entity.Address, entity.CreatedAt, entity.UpdatedAt))

		seq, err := repo.SeqByPeerId(peerId)
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		for item := range seq {
			assert.Equal(t, entity.ID, item.ID)
			assert.Equal(t, entity.AppNodeID, item.AppNodeID)
			assert.Equal(t, entity.Address, item.Address)
			assert.Equal(t, entity.CreatedAt, item.CreatedAt)
			assert.Equal(t, entity.UpdatedAt, item.UpdatedAt)
			count++
		}
		assert.Equal(t, 1, count)
	})

}
