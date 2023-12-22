package broadcast_test

import (
	"crypto/rand"
	"net"
	"pan/broadcast"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestModel ...
func TestModel(t *testing.T) {

	t.Run("Init", func(t *testing.T) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.Close()

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))
		mock.ExpectExec("CREATE TABLE `records`").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE INDEX `idx_seq`").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE INDEX `idx_addr`").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("CREATE INDEX `idx_records_deleted_at`").WillReturnResult(sqlmock.NewResult(0, 0))
		mock.MatchExpectationsInOrder(false)

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		model := broadcast.NewRepo(db)
		err = model.Init()
		if err != nil {
			t.Fatal(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("Save", func(t *testing.T) {

		token := make([]byte, 32)
		rand.Read(token)
		addr := []byte(net.JoinHostPort("127.0.0.1", "9000"))
		peerId := uuid.New()

		rd := new(broadcast.Record)
		rd.Seq = time.Now().Unix()
		rd.Token = token
		rd.Addr = addr
		rd.PeerId = peerId[:]
		rd.DeathTime = time.Now().Unix()

		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.Close()

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))
		mock.ExpectExec("INSERT INTO `records`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), rd.Seq, rd.Token, rd.Addr, rd.PeerId, rd.DeathTime).WillReturnResult(sqlmock.NewResult(1, 1))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo := broadcast.NewRepo(db)
		err = repo.Save(rd)
		if err != nil {
			t.Fatal(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("FindOneWithAddrAndSeq", func(t *testing.T) {

		token := make([]byte, 32)
		rand.Read(token)
		seq := time.Now().Unix()
		id := int64(123)
		addr := []byte(net.JoinHostPort("127.0.0.1", "9000"))

		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer mockDB.Close()

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))
		mock.ExpectQuery("SELECT (.+) FROM `records`").WillReturnRows(sqlmock.NewRows([]string{"id", "seq", "token", "addr"}).AddRow(id, seq, token, addr))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo := broadcast.NewRepo(db)
		rd, err := repo.FindOneWithAddrAndSeq(addr, seq)
		if err != nil {
			t.Fatal(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		assert.Equal(t, id, rd.ID, "ID should be same")
		assert.Equal(t, seq, rd.Seq, "Seq should be same")
		assert.Equal(t, token, rd.Token, "Token should be same")
		assert.Equal(t, addr, rd.Addr, "Addr should be same")

	})
}
