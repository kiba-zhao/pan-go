package repositories

import (
	"io/fs"
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type FileInfoIteration = func(fileInfo *models.TargetFile) error
type TargetFileRepository interface {
	FindOrCreateByTargetIDAndRelativePath(targetID uint, relativePath string) (models.TargetFile, error)
	Save(fileInfo models.TargetFile) error
	UpdateEachFileInfoByTargetID(targetID uint, iteration FileInfoIteration) error
}

type TargetFileRepositoryImpl struct {
	DB *gorm.DB
}

func NewTargetFileRepository(db *gorm.DB) TargetFileRepository {
	repo := new(TargetFileRepositoryImpl)
	repo.DB = db
	return repo
}

func (repo *TargetFileRepositoryImpl) FindOrCreateByTargetIDAndRelativePath(targetID uint, relativePath string) (models.TargetFile, error) {

	var model models.TargetFile
	model.TargetID = targetID
	model.RelativePath = relativePath
	results := repo.DB.FirstOrCreate(&model, model)
	return model, results.Error
}

func (repo *TargetFileRepositoryImpl) Save(fileInfo models.TargetFile) error {
	results := repo.DB.Save(&fileInfo)
	if results.Error == nil && results.RowsAffected != 1 {
		return gorm.ErrRecordNotFound
	}
	return results.Error
}

func (repo *TargetFileRepositoryImpl) UpdateEachFileInfoByTargetID(targetID uint, iteration FileInfoIteration) error {
	var model models.TargetFile
	model.TargetID = targetID
	rows, err := repo.DB.Model(&model).Where(&model).Rows()

	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {

		var fileInfo models.TargetFile
		err = repo.DB.ScanRows(rows, &fileInfo)
		if err != nil {
			break
		}

		modifyTime := fileInfo.ModifyTime
		err = iteration(&fileInfo)
		if err == fs.ErrNotExist {
			err = repo.DB.Delete(&fileInfo).Error
			continue
		} else if err == nil && modifyTime != fileInfo.ModifyTime {
			err = repo.DB.Save(&fileInfo).Error
		}

		if err != nil {
			break
		}
	}
	return err
}
