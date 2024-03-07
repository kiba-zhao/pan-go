package services

import (
	"pan/extfs/models"
	"pan/extfs/repositories"
)

type TargetService struct {
	TargetRepo repositories.TargetRepository
}

func (s *TargetService) Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error) {
	return s.TargetRepo.Search(conditions)
}
