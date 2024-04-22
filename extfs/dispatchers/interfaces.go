package dispatchers

import "pan/extfs/models"

type DispatchDone func(err error)

type TargetDispatcher interface {
	Scan(target models.Target, done DispatchDone) error
	Clean(target models.Target, done DispatchDone) error
}
