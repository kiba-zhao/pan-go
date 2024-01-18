package services

import (
	"pan/modules/extfs/models"
	"pan/modules/extfs/repositories"
)

type TagService struct {
	TagRepo repositories.TagRepository
}

// Find ...
func (s *TagService) Find(condition *models.TagFindCondition) (tags []models.Tag, err error) {
	tags, err = s.TagRepo.Find(condition)
	return
}
