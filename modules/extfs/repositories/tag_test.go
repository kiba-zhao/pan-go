package repositories_test

import (
	"database/sql"

	"fmt"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestTagRepository ...
func TestTagRepository(t *testing.T) {

	setup := func() (repo repositories.TagRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewTagRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("Find", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		id := uint(1)
		name := "tag name"
		owner := uuid.New().String()

		mock.ExpectQuery("SELECT (.+) FROM `tags` WHERE `tags`.`deleted_at` IS NULL").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner"}).AddRow(id, name, owner))

		tags, err := repo.Find(nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, id, tags[0].ID)
		assert.Equal(t, name, tags[0].Name)
		assert.Equal(t, owner, tags[0].Owner)
	})

	t.Run("Find with condition", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		id := uint(1)
		name := "tag name"
		owner := uuid.New().String()
		limit := 101
		offset := 1

		sql := fmt.Sprintf("SELECT (.+) FROM `tags` WHERE (.+) LIMIT %d OFFSET %d", limit, offset)
		mock.ExpectQuery(sql).WithArgs(owner, name+"%").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "owner"}).AddRow(id, name, owner))

		tags, err := repo.Find(&models.TagFindCondition{Name: name, Owner: owner, Limit: limit, Offset: offset})
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, tags, 1)
		assert.Equal(t, id, tags[0].ID)
		assert.Equal(t, name, tags[0].Name)
		assert.Equal(t, owner, tags[0].Owner)
	})
}
