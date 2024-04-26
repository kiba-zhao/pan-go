package services

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/repositories"
)

type TargetFileService struct {
	TargetFileRepo repositories.TargetFileRepository
}

func (s *TargetFileService) Search(conditions models.TargetFileSearchCondition) (total int64, items []models.TargetFile, err error) {
	total, items_, err := s.TargetFileRepo.Search(conditions, true)
	if err != nil {
		return
	}

	for _, item := range items_ {
		setTargetFileAvailable(&item, &item.Target)
		if conditions.Available == nil || *conditions.Available == item.Available {
			items = append(items, item)
		}
	}
	return
}

func (s *TargetFileService) Select(id uint64) (models.TargetFile, error) {
	targetFile, err := s.TargetFileRepo.Select(id, true)
	if err != nil {
		return targetFile, err
	}
	setTargetFileAvailable(&targetFile, &targetFile.Target)
	return targetFile, nil
}

func (s *TargetFileService) ScanByTarget(target models.Target) error {
	return s.TargetFileRepo.TraverseByTargetId(func(targetFile models.TargetFile) error {

		if target.HashCode != targetFile.TargetHashCode {
			return s.TargetFileRepo.Delete(targetFile)
		}

		var stat os.FileInfo
		isUpdated, err := exploreFile(&targetFile)
		if err == nil {
			stat, err = os.Stat(targetFile.FilePath)
		}

		if os.IsNotExist(err) {
			return s.TargetFileRepo.Delete(targetFile)
		}

		if err != nil {
			return err
		}

		if stat.ModTime() != targetFile.ModTime {
			targetFile.ModTime = stat.ModTime()
			isUpdated = true
		}

		if stat.Size() != targetFile.Size {
			targetFile.Size = stat.Size()
			isUpdated = true
		}

		if isUpdated {
			_, err = s.TargetFileRepo.Save(targetFile)
		}
		return err

	}, target.ID)
}

func (s *TargetFileService) CleanByTarget(target models.Target) error {
	return s.TargetFileRepo.DeleteByTargetId(target.ID)
}

func (s *TargetFileService) ScanFileByTarget(filepath string, target models.Target) error {
	hashCode := generateHashCodeByFilePath(filepath)
	targetFile, err := s.TargetFileRepo.SelectByFilePathAndTargetId(filepath, target.ID, hashCode, false)
	if err != nil && err != errors.ErrNotFound {
		return err
	}

	if err == errors.ErrNotFound {
		targetFile = models.TargetFile{}
		targetFile.FilePath = filepath
		targetFile.HashCode = hashCode
		targetFile.TargetID = target.ID
	}

	stat, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	targetFile.TargetHashCode = target.HashCode
	targetFile.ModTime = stat.ModTime()
	targetFile.Size = stat.Size()

	_, err = exploreFile(&targetFile)
	if err == nil {
		_, err = s.TargetFileRepo.Save(targetFile)
	}

	return err
}

func exploreFile(targetFile *models.TargetFile) (updated bool, err error) {
	file, err := os.Open(targetFile.FilePath)
	if err != nil {
		return
	}
	defer file.Close()

	mineBytes := make([]byte, 512)
	count, err := file.Read(mineBytes)
	if err != nil {
		return
	}

	mimeType := http.DetectContentType(mineBytes[:count])
	if mimeType != targetFile.MimeType {
		targetFile.MimeType = mimeType
		updated = true
	}

	hash := sha512.New()
	hash.Write(mineBytes[:count])
	if count == 512 {
		io.Copy(hash, file)
	}
	checkSum := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	if checkSum != targetFile.CheckSum {
		targetFile.CheckSum = checkSum
		updated = true
	}

	return
}

func setTargetFileAvailable(targetFile *models.TargetFile, target *models.Target) {
	targetFile.Available = *target.Enabled
	if targetFile.Available {
		setTargetAvailable(target)
		targetFile.Available = target.Available
	}
	if targetFile.Available {
		_, err := os.Stat(targetFile.FilePath)
		targetFile.Available = err == nil
	}

}
