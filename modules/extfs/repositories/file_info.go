package repositories

import (
	"io/fs"
	"pan/modules/extfs/models"

	"gorm.io/gorm"
)

type FileInfoIteration = func(fileInfo *models.FileInfo) error
type FileInfoRepository interface {
	FindOrCreateByTargetIDAndRelativePath(targetID uint, relativePath string) (models.FileInfo, error)
	Save(fileInfo models.FileInfo) error
	UpdateEachFileInfoByTargetID(targetID uint, iteration FileInfoIteration) error
}

type fileInfoRepositoryImpl struct {
	DB *gorm.DB
}

func NewFileInfoRepository(db *gorm.DB) FileInfoRepository {
	repo := new(fileInfoRepositoryImpl)
	repo.DB = db
	return repo
}

func (repo *fileInfoRepositoryImpl) FindOrCreateByTargetIDAndRelativePath(targetID uint, relativePath string) (models.FileInfo, error) {

	var model models.FileInfo
	model.TargetID = targetID
	model.RelativePath = relativePath
	results := repo.DB.FirstOrCreate(&model, model)
	return model, results.Error
}

func (repo *fileInfoRepositoryImpl) Save(fileInfo models.FileInfo) error {
	results := repo.DB.Save(&fileInfo)
	return results.Error
}

func (repo *fileInfoRepositoryImpl) UpdateEachFileInfoByTargetID(targetID uint, iteration FileInfoIteration) error {
	rows, err := repo.DB.Where("target_id = ? ", targetID).Rows()
	defer rows.Close()
	if err != nil {
		return err
	}

	for rows.Next() {

		var fileInfo models.FileInfo
		err = repo.DB.ScanRows(rows, &fileInfo)
		if err != nil {
			break
		}

		modifyTime := fileInfo.ModifyTime
		err := iteration(&fileInfo)
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
