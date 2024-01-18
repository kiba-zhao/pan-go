package repositories

import "pan/modules/extfs/models"

type ExtFSRepository interface {
	GetLatestOne() (models.ExtFS, error)
}
