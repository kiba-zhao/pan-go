package repositories_test

import (
	"database/sql"
	"io/fs"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFileInfoRepository(t *testing.T) {

	setup := func() (repo repositories.FileInfoRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = repositories.NewFileInfoRepository(db)
		return
	}

	teardown := func(mockDB *sql.DB) {
		mockDB.Close()
	}

	t.Run("FindOrCreateByTargetIDAndRelativePath", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var fileInfoRow models.FileInfo
		fileInfoRow.ID = 1
		fileInfoRow.TargetID = 1
		fileInfoRow.RelativePath = "path"
		fileInfoRow.Hash = []byte("hash")
		fileInfoRow.ModifyTime = time.Now()
		fileInfoRow.Name = "name"
		fileInfoRow.Size = 100

		mock.ExpectQuery("SELECT  (.+) FROM `file_infos` WHERE `file_infos`.`target_id` = \\? AND `file_infos`.`relative_path` = \\? AND `file_infos`.`deleted_at` IS NULL ORDER BY `file_infos`.`id` LIMIT 1").WithArgs(fileInfoRow.TargetID, fileInfoRow.RelativePath).WillReturnRows(sqlmock.NewRows([]string{"id", "target_id", "relative_path", "hash", "modify_time", "name", "size"}).AddRow(fileInfoRow.ID, fileInfoRow.TargetID, fileInfoRow.RelativePath, fileInfoRow.Hash, fileInfoRow.ModifyTime, fileInfoRow.Name, fileInfoRow.Size))

		row, err := repo.FindOrCreateByTargetIDAndRelativePath(fileInfoRow.TargetID, fileInfoRow.RelativePath)

		assert.Nil(t, err)
		assert.Equal(t, fileInfoRow.ID, row.ID)
		assert.Equal(t, fileInfoRow.TargetID, row.TargetID)
		assert.Equal(t, fileInfoRow.RelativePath, row.RelativePath)
		assert.Equal(t, fileInfoRow.Hash, row.Hash)
		assert.Equal(t, fileInfoRow.ModifyTime, row.ModifyTime)
		assert.Equal(t, fileInfoRow.Name, row.Name)
		assert.Equal(t, fileInfoRow.Size, row.Size)

	})

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var fileInfoRow models.FileInfo
		fileInfoRow.TargetID = 1
		fileInfoRow.RelativePath = "path"
		fileInfoRow.Hash = []byte("hash")
		fileInfoRow.ModifyTime = time.Now()
		fileInfoRow.Name = "name"
		fileInfoRow.Size = 100

		mock.ExpectExec("INSERT INTO `file_infos`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), fileInfoRow.TargetID, fileInfoRow.Name, fileInfoRow.Size, fileInfoRow.ModifyTime, fileInfoRow.RelativePath, fileInfoRow.Hash).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Save(fileInfoRow)
		assert.Nil(t, err)

		fileInfoRow.ID = 1
		mock.ExpectExec("UPDATE `file_infos` SET (.+) WHERE `file_infos`.`deleted_at` IS NULL AND `id` = \\?").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), fileInfoRow.TargetID, fileInfoRow.Name, fileInfoRow.Size, fileInfoRow.ModifyTime, fileInfoRow.RelativePath, fileInfoRow.Hash, fileInfoRow.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.Save(fileInfoRow)
		assert.Nil(t, err)
	})

	t.Run("UpdateEachFileInfoByTargetID", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer teardown(mockDB)

		var fileInfoRow models.FileInfo
		fileInfoRow.ID = 1
		fileInfoRow.TargetID = 1
		fileInfoRow.RelativePath = "path"
		fileInfoRow.Hash = []byte("hash")
		fileInfoRow.ModifyTime = time.Now().Add(-24 * time.Hour)
		fileInfoRow.Name = "name"
		fileInfoRow.Size = 100
		modifyTime := time.Now()

		mock.ExpectQuery("SELECT  (.+) FROM `file_infos` WHERE `file_infos`.`target_id` = \\? AND `file_infos`.`deleted_at` IS NULL").WithArgs(fileInfoRow.TargetID).WillReturnRows(sqlmock.NewRows([]string{"id", "target_id", "relative_path", "hash", "modify_time", "name", "size"}).AddRow(fileInfoRow.ID, fileInfoRow.TargetID, fileInfoRow.RelativePath, fileInfoRow.Hash, fileInfoRow.ModifyTime, fileInfoRow.Name, fileInfoRow.Size))
		mock.ExpectExec("UPDATE `file_infos` SET (.+) WHERE `file_infos`.`deleted_at` IS NULL AND `id` = \\?").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), fileInfoRow.TargetID, fileInfoRow.Name, fileInfoRow.Size, modifyTime, fileInfoRow.RelativePath, fileInfoRow.Hash, fileInfoRow.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.UpdateEachFileInfoByTargetID(fileInfoRow.TargetID, func(fileInfo *models.FileInfo) error {
			assert.Equal(t, fileInfoRow.ID, fileInfo.ID)
			assert.Equal(t, fileInfoRow.TargetID, fileInfo.TargetID)
			assert.Equal(t, fileInfoRow.RelativePath, fileInfo.RelativePath)
			assert.Equal(t, fileInfoRow.Hash, fileInfo.Hash)
			assert.Equal(t, fileInfoRow.ModifyTime, fileInfo.ModifyTime)
			assert.Equal(t, fileInfoRow.Name, fileInfo.Name)
			assert.Equal(t, fileInfoRow.Size, fileInfo.Size)
			fileInfo.ModifyTime = modifyTime
			return nil
		})
		assert.Nil(t, err)

		mock.ExpectQuery("SELECT  (.+) FROM `file_infos` WHERE `file_infos`.`target_id` = \\? AND `file_infos`.`deleted_at` IS NULL").WithArgs(fileInfoRow.TargetID).WillReturnRows(sqlmock.NewRows([]string{"id", "target_id", "relative_path", "hash", "modify_time", "name", "size"}).AddRow(fileInfoRow.ID, fileInfoRow.TargetID, fileInfoRow.RelativePath, fileInfoRow.Hash, fileInfoRow.ModifyTime, fileInfoRow.Name, fileInfoRow.Size))
		mock.ExpectExec("UPDATE `file_infos` SET `deleted_at`=\\? WHERE `file_infos`.`id` = \\? AND `file_infos`.`deleted_at` IS NULL").WithArgs(sqlmock.AnyArg(), fileInfoRow.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.UpdateEachFileInfoByTargetID(fileInfoRow.TargetID, func(fileInfo *models.FileInfo) error {
			assert.Equal(t, fileInfoRow.ID, fileInfo.ID)
			assert.Equal(t, fileInfoRow.TargetID, fileInfo.TargetID)
			assert.Equal(t, fileInfoRow.RelativePath, fileInfo.RelativePath)
			assert.Equal(t, fileInfoRow.Hash, fileInfo.Hash)
			assert.Equal(t, fileInfoRow.ModifyTime, fileInfo.ModifyTime)
			assert.Equal(t, fileInfoRow.Name, fileInfo.Name)
			assert.Equal(t, fileInfoRow.Size, fileInfo.Size)
			return fs.ErrNotExist
		})
		assert.Nil(t, err)
	})

}
