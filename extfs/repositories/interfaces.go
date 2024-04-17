package repositories

import "pan/extfs/models"

type TargetRepository interface {
	Search(conditions models.TargetSearchCondition) (total int64, items []models.Target, err error)
	Save(target models.Target, withVersion bool) (models.Target, error)
	Select(id uint, version *uint8) (models.Target, error)
	Delete(target models.Target) error
}

type TargetFileTraverse = func(targetFile models.TargetFile) error
type TargetFileRepository interface {
	TraverseByTargetId(f TargetFileTraverse, targetId uint) error
	Save(targetFile models.TargetFile) (models.TargetFile, error)
	Delete(targetFile models.TargetFile) error
	DeleteByTargetId(targetId uint) error
	Select(id uint64) (models.TargetFile, error)
	SelectByFilePathAndTargetId(filepath string, targetId uint, hashCode string) (models.TargetFile, error)
}
