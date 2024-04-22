package impl_test

import (
	"database/sql"
	"errors"
	"pan/extfs/models"
	"pan/extfs/repositories"
	"pan/extfs/repositories/impl"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTargetFile(t *testing.T) {

	setup := func() (repo repositories.TargetFileRepository, mockDB *sql.DB, mock sqlmock.Sqlmock) {
		mockDB, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{""}).AddRow("3.8.10"))

		db, err := gorm.Open(sqlite.Dialector{Conn: mockDB}, &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			t.Fatal(err)
		}

		repo = &impl.TargetFileRepository{DB: db}
		return

	}

	t.Run("Save", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetFile := models.TargetFile{}
		targetFile.TargetID = 1
		targetFile.TargetHashCode = "target hash code"
		targetFile.FilePath = "file path"
		targetFile.HashCode = "hash code"
		targetFile.MimeType = "mime type"
		targetFile.CheckSum = "check sum"
		targetFile.ModTime = time.Now()
		targetFile.Size = 1

		mock.ExpectExec("INSERT INTO `target_files`").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetFile.TargetID, targetFile.TargetHashCode, targetFile.HashCode, targetFile.FilePath,
			targetFile.MimeType, targetFile.Size, targetFile.ModTime, targetFile.CheckSum,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Save(targetFile)
		assert.Nil(t, err)
		assert.Greater(t, result.ID, uint64(0))
		assert.Equal(t, targetFile.TargetID, result.TargetID)
		assert.Equal(t, targetFile.TargetHashCode, result.TargetHashCode)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.HashCode, result.HashCode)
		assert.Equal(t, targetFile.MimeType, result.MimeType)
		assert.Equal(t, targetFile.CheckSum, result.CheckSum)
		assert.Equal(t, targetFile.ModTime, result.ModTime)
		assert.Equal(t, targetFile.Size, result.Size)

		targetFile.ID = 2
		mock.ExpectExec("UPDATE `target_files`").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), targetFile.TargetID, targetFile.TargetHashCode, targetFile.HashCode, targetFile.FilePath,
			targetFile.MimeType, targetFile.Size, targetFile.ModTime, targetFile.CheckSum, targetFile.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		result, err = repo.Save(targetFile)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, targetFile.TargetID, result.TargetID)
		assert.Equal(t, targetFile.TargetHashCode, result.TargetHashCode)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.HashCode, result.HashCode)
		assert.Equal(t, targetFile.MimeType, result.MimeType)
		assert.Equal(t, targetFile.CheckSum, result.CheckSum)
		assert.Equal(t, targetFile.ModTime, result.ModTime)
		assert.Equal(t, targetFile.Size, result.Size)

	})

	t.Run("Delete", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetFile := models.TargetFile{}
		targetFile.ID = 1

		mock.ExpectExec("UPDATE `target_files`").WithArgs(sqlmock.AnyArg(), targetFile.ID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(targetFile)
		assert.Nil(t, err)
	})

	t.Run("DeleteByTargetId", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetID := uint(1)

		mock.ExpectExec("UPDATE `target_files`").WithArgs(sqlmock.AnyArg(), targetID).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteByTargetId(targetID)
		assert.Nil(t, err)
	})

	t.Run("Select", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetFile := models.TargetFile{}
		targetFile.ID = 1
		targetFile.TargetID = 1
		targetFile.TargetHashCode = "target hash code"
		targetFile.FilePath = "file path"
		targetFile.HashCode = "hash code"
		targetFile.MimeType = "mime type"
		targetFile.CheckSum = "check sum"
		targetFile.ModTime = time.Now()
		targetFile.Size = 1

		mock.ExpectQuery("SELECT (.+) FROM `target_files`").WithArgs(targetFile.ID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "target_id", "target_hash_code", "file_path", "hash_code", "mime_type", "check_sum", "mod_time", "size"}).
				AddRow(targetFile.ID, targetFile.TargetID, targetFile.TargetHashCode, targetFile.FilePath, targetFile.HashCode, targetFile.MimeType, targetFile.CheckSum, targetFile.ModTime, targetFile.Size),
		)

		result, err := repo.Select(targetFile.ID)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, targetFile.TargetID, result.TargetID)
		assert.Equal(t, targetFile.TargetHashCode, result.TargetHashCode)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.HashCode, result.HashCode)
		assert.Equal(t, targetFile.MimeType, result.MimeType)
		assert.Equal(t, targetFile.CheckSum, result.CheckSum)
		assert.Equal(t, targetFile.ModTime, result.ModTime)
		assert.Equal(t, targetFile.Size, result.Size)
	})

	t.Run("SelectByFilePathAndTargetId", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetFile := models.TargetFile{}
		targetFile.ID = 1
		targetFile.TargetID = 1
		targetFile.TargetHashCode = "target hash code"
		targetFile.FilePath = "file path"
		targetFile.HashCode = "hash code"
		targetFile.MimeType = "mime type"
		targetFile.CheckSum = "check sum"
		targetFile.ModTime = time.Now()
		targetFile.Size = 1

		mock.ExpectQuery("SELECT (.+) FROM `target_files`").WithArgs(targetFile.TargetID, targetFile.HashCode, targetFile.FilePath).WillReturnRows(
			sqlmock.NewRows([]string{"id", "target_id", "target_hash_code", "file_path", "hash_code", "mime_type", "check_sum", "mod_time", "size"}).
				AddRow(targetFile.ID, targetFile.TargetID, targetFile.TargetHashCode, targetFile.FilePath, targetFile.HashCode, targetFile.MimeType, targetFile.CheckSum, targetFile.ModTime, targetFile.Size),
		)

		result, err := repo.SelectByFilePathAndTargetId(targetFile.FilePath, targetFile.TargetID, targetFile.HashCode)
		assert.Nil(t, err)
		assert.Equal(t, targetFile.ID, result.ID)
		assert.Equal(t, targetFile.TargetID, result.TargetID)
		assert.Equal(t, targetFile.TargetHashCode, result.TargetHashCode)
		assert.Equal(t, targetFile.FilePath, result.FilePath)
		assert.Equal(t, targetFile.HashCode, result.HashCode)
		assert.Equal(t, targetFile.MimeType, result.MimeType)
		assert.Equal(t, targetFile.CheckSum, result.CheckSum)
		assert.Equal(t, targetFile.ModTime, result.ModTime)
		assert.Equal(t, targetFile.Size, result.Size)
	})

	t.Run("TraverseByTargetId", func(t *testing.T) {
		repo, mockDB, mock := setup()
		defer mockDB.Close()

		targetFile := models.TargetFile{}
		targetFile.ID = 1
		targetFile.TargetID = 1
		targetFile.TargetHashCode = "target hash code"
		targetFile.FilePath = "file path"
		targetFile.HashCode = "hash code"
		targetFile.MimeType = "mime type"
		targetFile.CheckSum = "check sum"
		targetFile.ModTime = time.Now()
		targetFile.Size = 1

		mock.ExpectQuery("SELECT (.+) FROM `target_files`").WithArgs(targetFile.TargetID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "target_id", "target_hash_code", "file_path", "hash_code", "mime_type", "check_sum", "mod_time", "size"}).
				AddRow(targetFile.ID, targetFile.TargetID, targetFile.TargetHashCode, targetFile.FilePath, targetFile.HashCode, targetFile.MimeType, targetFile.CheckSum, targetFile.ModTime, targetFile.Size),
		)

		err := repo.TraverseByTargetId(func(row models.TargetFile) error {
			assert.Equal(t, targetFile.ID, row.ID)
			assert.Equal(t, targetFile.TargetID, row.TargetID)
			assert.Equal(t, targetFile.TargetHashCode, row.TargetHashCode)
			assert.Equal(t, targetFile.FilePath, row.FilePath)
			assert.Equal(t, targetFile.HashCode, row.HashCode)
			assert.Equal(t, targetFile.MimeType, row.MimeType)
			assert.Equal(t, targetFile.CheckSum, row.CheckSum)
			assert.Equal(t, targetFile.ModTime, row.ModTime)
			assert.Equal(t, targetFile.Size, row.Size)
			return nil
		}, targetFile.TargetID)
		assert.Nil(t, err)

		mock.ExpectQuery("SELECT (.+) FROM `target_files`").WithArgs(targetFile.TargetID).WillReturnRows(
			sqlmock.NewRows([]string{"id", "target_id", "target_hash_code", "file_path", "hash_code", "mime_type", "check_sum", "mod_time", "size"}).
				AddRow(targetFile.ID, targetFile.TargetID, targetFile.TargetHashCode, targetFile.FilePath, targetFile.HashCode, targetFile.MimeType, targetFile.CheckSum, targetFile.ModTime, targetFile.Size),
		)

		testErr := errors.New("test error")
		err = repo.TraverseByTargetId(func(row models.TargetFile) error {
			return testErr
		}, targetFile.TargetID)
		assert.Equal(t, testErr, err)

	})
}
