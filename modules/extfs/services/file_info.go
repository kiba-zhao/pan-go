package services

import (
	"bytes"
	"crypto/sha512"
	"io"
	"io/fs"
	"os"
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

const CHUNK_SIZE = 64 * 1024

type FileInfoService struct {
	FileInfoRepo repositories.FileInfoRepository
}

func (s *FileInfoService) UpgradeFileInfoForTarget(target models.Target) error {

	dirFS := os.DirFS(target.FilePath)
	return s.FileInfoRepo.UpdateEachFileInfoByTargetID(target.ID, func(fileInfo *models.FileInfo) (err error) {

		file, err := dirFS.Open(fileInfo.RelativePath)
		if err != nil {
			switch err.(type) {
			case *fs.PathError:
				err = fs.ErrNotExist
			}
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil || fileInfo.ModifyTime == stat.ModTime() {
			return
		}
		fileInfo.ModifyTime = stat.ModTime()

		sig, err := s.GenerateFileSignature(file)
		if err == nil && !s.EqualFileInfo(*fileInfo, stat.Name(), stat.Size(), sig) {
			fileInfo.Name = stat.Name()
			fileInfo.Size = stat.Size()
			fileInfo.Hash = sig
		}

		return
	})
}

func (s *FileInfoService) ScanFileInfoForTarget(target models.Target) error {
	fileInfo, err := s.FileInfoRepo.FindOrCreateByTargetIDAndRelativePath(target.ID, target.FilePath)
	if err != nil || fileInfo.ModifyTime == target.ModifyTime {
		return err
	}
	fileInfo.ModifyTime = target.ModifyTime

	reader, err := os.Open(target.FilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	sig, err := s.GenerateFileSignature(reader)
	if err != nil {
		return err
	}

	if !s.EqualFileInfo(fileInfo, target.Name, target.Size, sig) {
		fileInfo.Name = target.Name
		fileInfo.Size = target.Size
		fileInfo.Hash = sig
	}

	err = s.FileInfoRepo.Save(fileInfo)

	return err
}

func (s *FileInfoService) ScanFileInfosForDirectory(target models.Target) (models.FileInfoTotal, error) {
	dirFS := os.DirFS(target.FilePath)
	var fileInfoTotal models.FileInfoTotal

	err := fs.WalkDir(dirFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		stat, err := d.Info()
		if err != nil {
			return err
		}

		fileInfo, err := s.FileInfoRepo.FindOrCreateByTargetIDAndRelativePath(target.ID, path)
		if err != nil || stat.ModTime() == fileInfo.ModifyTime {
			return err
		}
		fileInfo.ModifyTime = stat.ModTime()

		reader, err := dirFS.Open(path)
		if err != nil {
			return err
		}
		defer reader.Close()

		sig, err := s.GenerateFileSignature(reader)
		if err != nil {
			return err
		}

		if !s.EqualFileInfo(fileInfo, stat.Name(), stat.Size(), sig) {
			fileInfo.Name = stat.Name()
			fileInfo.Size = stat.Size()
			fileInfo.Hash = sig
		}

		err = s.FileInfoRepo.Save(fileInfo)

		if err == nil {
			fileInfoTotal.Size += fileInfo.Size
			fileInfoTotal.Total++
		}

		return err
	})

	return fileInfoTotal, err
}

func (s *FileInfoService) GenerateFileSignature(reader io.Reader) ([]byte, error) {

	hash := sha512.New()
	var err error
	for {
		_, copyErr := io.CopyN(hash, reader, CHUNK_SIZE)
		if copyErr != nil {
			if copyErr != io.EOF {
				err = copyErr
			}
			break
		}
	}

	return hash.Sum(nil), err
}

func (s *FileInfoService) EqualFileInfo(fileInfo models.FileInfo, name string, size int64, hash []byte) bool {
	if fileInfo.Name != name {
		return false
	}
	if fileInfo.Size != size {
		return false
	}
	if !bytes.Equal(fileInfo.Hash, hash) {
		return false
	}
	return true
}
