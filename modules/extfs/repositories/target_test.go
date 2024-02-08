package repositories_test

import (
	"database/sql"
	"fmt"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTarget(t *testing.T) {

	setup := func() (repo repositories.TargetRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewTargetRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("FindAllWithEnabled", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var targetRow models.Target
		targetRow.ID = 1
		targetRow.Enabled = true
		targetRow.Name = "target name"
		targetRow.FilePath = "target file path"
		targetRow.Size = 100
		targetRow.Total = 10
		targetRow.ModifyTime = time.Now()
		targetRow.CreatedAt = time.Now()
		targetRow.UpdatedAt = time.Now()

		mock.ExpectQuery("SELECT  (.+) FROM `targets` WHERE `targets`.`enabled` = \\? AND `targets`.`deleted_at` IS NULL").WithArgs(true).WillReturnRows(sqlmock.NewRows([]string{"id", "enabled", "name", "file_path", "size", "total", "modify_time", "created_at", "updated_at"}).
			AddRow(targetRow.ID, targetRow.Enabled, targetRow.Name, targetRow.FilePath, targetRow.Size, targetRow.Total, targetRow.ModifyTime, targetRow.CreatedAt, targetRow.UpdatedAt))

		rows, err := repo.FindAllWithEnabled()

		assert.NoError(t, err)
		assert.Equal(t, targetRow.ID, rows[0].ID)
		assert.Equal(t, targetRow.Name, rows[0].Name)
		assert.Equal(t, targetRow.FilePath, rows[0].FilePath)
		assert.Equal(t, targetRow.Size, rows[0].Size)
		assert.Equal(t, targetRow.Total, rows[0].Total)
		assert.Equal(t, targetRow.Enabled, rows[0].Enabled)
		assert.Equal(t, targetRow.CreatedAt, rows[0].CreatedAt)
		assert.Equal(t, targetRow.UpdatedAt, rows[0].UpdatedAt)
	})

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		targetRow := models.Target{}
		targetRow.Enabled = true
		targetRow.Name = "target name"
		targetRow.FilePath = "target file path"
		targetRow.Size = 100
		targetRow.Total = 10
		targetRow.ModifyTime = time.Now()

		mock.ExpectExec("INSERT INTO `targets`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetRow.FilePath, targetRow.Name, targetRow.Size, targetRow.ModifyTime, targetRow.Enabled, targetRow.Total).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(targetRow)

		assert.NoError(t, err)

		targetRow.ID = 1
		mock.ExpectExec("UPDATE `targets` SET (.+) WHERE `targets`.`deleted_at` IS NULL AND `id` = \\?").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetRow.FilePath, targetRow.Name, targetRow.Size, targetRow.ModifyTime, targetRow.Enabled, targetRow.Total, targetRow.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Save(targetRow)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var condition models.TargetSearchCondition

		condition.Limit = 321
		condition.Offset = 32
		condition.Keyword = " key1, key2 ,key 3"

		total := int64(123)
		mock.ExpectQuery("SELECT  (.+) FROM `targets` WHERE (.+) AND `targets`.`deleted_at` IS NULL").WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(total))

		var targetRow models.Target
		targetRow.ID = 1
		targetRow.Enabled = true
		targetRow.Name = "target name"
		targetRow.FilePath = "target file path"
		targetRow.Size = 100
		targetRow.Total = 10
		targetRow.ModifyTime = time.Now()
		targetRow.CreatedAt = time.Now()
		targetRow.UpdatedAt = time.Now()

		sql := fmt.Sprintf("SELECT (.+) FROM `targets` WHERE (.+) LIMIT %d OFFSET %d", condition.Limit, condition.Offset)
		mock.ExpectQuery(sql).WillReturnRows(sqlmock.NewRows([]string{"id", "enabled", "name", "file_path", "size", "total", "modify_time", "created_at", "updated_at"}).
			AddRow(targetRow.ID, targetRow.Enabled, targetRow.Name, targetRow.FilePath, targetRow.Size, targetRow.Total, targetRow.ModifyTime, targetRow.CreatedAt, targetRow.UpdatedAt))

		rTotal, rTargets, err := repo.Search(&condition)

		assert.NoError(t, err)
		assert.Equal(t, total, rTotal)
		assert.Equal(t, targetRow.ID, rTargets[0].ID)
		assert.Equal(t, targetRow.Name, rTargets[0].Name)
		assert.Equal(t, targetRow.FilePath, rTargets[0].FilePath)
		assert.Equal(t, targetRow.Size, rTargets[0].Size)
		assert.Equal(t, targetRow.Total, rTargets[0].Total)
		assert.Equal(t, targetRow.Enabled, rTargets[0].Enabled)
		assert.Equal(t, targetRow.CreatedAt, rTargets[0].CreatedAt)
		assert.Equal(t, targetRow.UpdatedAt, rTargets[0].UpdatedAt)
	})

}
