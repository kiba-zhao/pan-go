package services

import (
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

type KeywordService struct {
	KeywordRepo repositories.KeywordRepository
}

// Find ...
func (s *KeywordService) Find(condition *models.KeywordFindCondition) (keywords []models.Keyword, err error) {
	keywords, err = s.KeywordRepo.Find(condition)
	return
}
