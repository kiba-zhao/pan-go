package impl

import (
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/repositories"

	"gorm.io/gorm"
)

type TargetFileRepository struct {
	DB *gorm.DB
}

func (repo *TargetFileRepository) Save(targetFile models.TargetFile) (models.TargetFile, error) {
	results := repo.DB.Save(&targetFile)
	if results.Error == nil && results.RowsAffected != 1 {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) Delete(targetFile models.TargetFile) error {
	results := repo.DB.Delete(&targetFile)
	if results.Error == nil && results.RowsAffected != 1 {
		return errors.ErrNotFound
	}
	return results.Error
}

func (repo *TargetFileRepository) DeleteByTargetId(targetId uint) error {
	var targetFile models.TargetFile
	targetFile.TargetID = targetId

	results := repo.DB.Where(&targetFile).Delete(&targetFile)
	return results.Error
}

func (repo *TargetFileRepository) Select(id uint64) (models.TargetFile, error) {
	var targetFile models.TargetFile
	results := repo.DB.Take(&targetFile, id)
	if results.Error == gorm.ErrRecordNotFound {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) SelectByFilePathAndTargetId(filepath string, targetId uint, hashCode string) (models.TargetFile, error) {
	var targetFile models.TargetFile
	targetFile.FilePath = filepath
	targetFile.TargetID = targetId
	targetFile.HashCode = hashCode

	results := repo.DB.Where(&targetFile).Take(&targetFile)
	if results.Error == gorm.ErrRecordNotFound {
		return targetFile, errors.ErrNotFound
	}
	return targetFile, results.Error
}

func (repo *TargetFileRepository) TraverseByTargetId(f repositories.TargetFileTraverse, targetId uint) error {
	var targetFile models.TargetFile
	targetFile.TargetID = targetId

	rows, err := repo.DB.Model(&targetFile).Where(&targetFile).Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var targetFile_ models.TargetFile
		err = repo.DB.ScanRows(rows, &targetFile_)
		if err != nil {
			break
		}
		err = f(targetFile_)
		if err != nil {
			break
		}
	}

	return err
}
