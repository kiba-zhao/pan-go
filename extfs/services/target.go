package services

import (
	"crypto/sha1"
	"encoding/base64"
	"io/fs"
	"os"
	"pan/extfs/dispatchers"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/repositories"
)

type TargetService struct {
	TargetRepo        repositories.TargetRepository
	TargetDispatcher  dispatchers.TargetDispatcher
	TargetFileService *TargetFileService
}

func (s *TargetService) Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error) {
	return s.TargetRepo.Search(conditions)
}

func (s *TargetService) Create(fields models.TargetFields) (models.Target, error) {
	_, err := os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return models.Target{}, err
	}
	available := err == nil
	version := uint8(1)

	var target models.Target
	target.Name = fields.Name
	target.HashCode = generateHashCodeByFilePath(fields.FilePath)
	target.FilePath = fields.FilePath
	target.Enabled = fields.Enabled
	target.Available = &available
	target.Version = &version

	target_, err := s.TargetRepo.Save(target, false)
	if err == nil {
		err = s.TargetDispatcher.Scan(target_)
	}
	return target_, err
}

func (s *TargetService) Update(fields models.TargetFields, id uint, opts models.TargetQueryOptions) (models.Target, error) {

	target, err := s.TargetRepo.Select(id, opts.Version)
	if err != nil {
		return target, err
	}

	_, err = os.Stat(fields.FilePath)
	if err != nil && !os.IsNotExist(err) {
		return target, err
	}
	available := err == nil
	err = nil

	if target.FilePath != fields.FilePath {
		target.HashCode = generateHashCodeByFilePath(fields.FilePath)
	}

	target.Name = fields.Name
	target.FilePath = fields.FilePath
	target.Enabled = fields.Enabled
	target.Available = &available

	target, err = s.TargetRepo.Save(target, true)
	if err == errors.ErrNotFound {
		return target, errors.ErrConflict
	}
	if err == nil {
		err = s.TargetDispatcher.Scan(target)
	}
	return target, err
}

func (s *TargetService) Select(id uint, opts models.TargetQueryOptions) (models.Target, error) {
	return s.TargetRepo.Select(id, opts.Version)
}

func (s *TargetService) Delete(id uint, opts models.TargetQueryOptions) error {
	target, err := s.TargetRepo.Select(id, opts.Version)
	if err != nil {
		return err
	}

	err = s.TargetRepo.Delete(target)
	if err == errors.ErrNotFound {
		return errors.ErrConflict
	}
	if err == nil {
		err = s.TargetDispatcher.Clean(target)
	}
	return err
}

func (s *TargetService) Scan(id uint) error {
	// TODO: Scenes that have been cleaned up
	target, err := s.TargetRepo.Select(id, nil)
	if err != nil {
		return err
	}

	if target.Available == nil || !*target.Available || target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	stat, err := os.Stat(target.FilePath)
	if err != nil {
		return err
	}

	err = s.TargetFileService.ScanByTarget(target)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return s.TargetFileService.ScanFileByTarget(target.FilePath, target)
	}

	return fs.WalkDir(os.DirFS(target.FilePath), target.FilePath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || err != nil {
			return err
		}
		return s.TargetFileService.ScanFileByTarget(path, target)
	})
}

func (s *TargetService) Clean(id uint) error {
	target, err := s.TargetRepo.Select(id, nil)
	if err != nil {
		return err
	}
	if !target.DeletedAt.Valid {
		return errors.ErrConflict
	}
	return s.TargetFileService.CleanByTarget(target)
}

func generateHashCodeByFilePath(filepath string) string {
	hash := sha1.Sum([]byte(filepath))
	return base64.StdEncoding.EncodeToString(hash[:])
}
