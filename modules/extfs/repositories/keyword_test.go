package repositories_test

import (
	"database/sql"

	"fmt"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestKeywordRepository ...
func TestKeywordRepository(t *testing.T) {

	setup := func() (repo repositories.KeywordRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewKeywordRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("Find", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		id := uint(1)
		name := "keyword name"

		mock.ExpectQuery("SELECT (.+) FROM `keywords` WHERE `keywords`.`deleted_at` IS NULL").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, name))

		keywords, err := repo.Find(nil)
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, keywords, 1)
		assert.Equal(t, id, keywords[0].ID)
		assert.Equal(t, name, keywords[0].Name)
	})

	t.Run("Find with condition", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		id := uint(1)
		name := "keyword name"
		limit := 101
		offset := 1

		sql := fmt.Sprintf("SELECT (.+) FROM `keywords` WHERE (.+) LIMIT %d OFFSET %d", limit, offset)
		mock.ExpectQuery(sql).WithArgs(name + "%").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(id, name))

		keywords, err := repo.Find(&models.KeywordFindCondition{Keyword: name, Limit: limit, Offset: offset})
		if err != nil {
			t.Fatal(err)
		}

		assert.Len(t, keywords, 1)
		assert.Equal(t, id, keywords[0].ID)
		assert.Equal(t, name, keywords[0].Name)
	})
}
