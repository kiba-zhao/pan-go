package dispatchers

import "pan/extfs/models"

type TargetDispatcher interface {
	Scan(target models.Target) error
	Clean(target models.Target) error
}
