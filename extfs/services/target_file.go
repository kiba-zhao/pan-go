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

func (s *TargetFileService) ScanByTarget(target models.Target) error {
	return s.TargetFileRepo.TraverseByTargetId(func(targetFile models.TargetFile) error {

		if target.HashCode != targetFile.TargetHashCode {
			return s.TargetFileRepo.Delete(targetFile)
		}

		isUpdated, err := exploreFile(&targetFile)
		if err != nil {
			return err
		}

		stat, err := os.Stat(targetFile.FilePath)
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
	targetFile, err := s.TargetFileRepo.SelectByFilePathAndTargetId(filepath, target.ID, hashCode)
	if err != nil && err != errors.ErrNotFound {
		return err
	}

	stat, err := os.Stat(filepath)
	if err != nil {
		return err
	}

	if err == errors.ErrNotFound {
		targetFile = models.TargetFile{}
		targetFile.FilePath = filepath
		targetFile.HashCode = hashCode
		targetFile.TargetID = target.ID
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
