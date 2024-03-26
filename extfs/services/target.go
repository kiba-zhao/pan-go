package services

import (
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/repositories"
)

type TargetService struct {
	TargetRepo repositories.TargetRepository
}

func (s *TargetService) Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error) {
	return s.TargetRepo.Search(conditions)
}

func (s *TargetService) Create(fields models.TargetFields) (models.Target, error) {
	var target models.Target
	target.Name = fields.Name
	target.FilePath = fields.FilePath
	target.Enabled = fields.Enabled
	return s.TargetRepo.Save(target, false)
}

func (s *TargetService) Update(fields models.TargetFields, id uint, opts models.TargetQueryOptions) (models.Target, error) {

	target, err := s.TargetRepo.Select(id, opts.Version)
	if err != nil {
		return target, err
	}

	target.Name = fields.Name
	target.FilePath = fields.FilePath
	target.Enabled = fields.Enabled

	target, err = s.TargetRepo.Save(target, true)
	if err == errors.ErrNotFound {
		return target, errors.ErrConflict
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
	return err
}
