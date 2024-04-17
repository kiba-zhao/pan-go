package impl

import (
	"cmp"
	"pan/cache"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
	"sync"
)

type TargetDispatcherItem struct {
	id         uint
	pending    *models.Target
	processing *models.Target
	sync.Mutex
}

func (t *TargetDispatcherItem) HashCode() uint {
	return t.id
}

type TargetDispatcherBucket = *cache.Bucket[uint, *TargetDispatcherItem]

func NewTargetDispatcherBucket() TargetDispatcherBucket {
	return cache.NewBucket[uint, *TargetDispatcherItem](cmp.Compare[uint])
}

type TargetDispatcher struct {
	TargetService *services.TargetService
	Bucket        TargetDispatcherBucket
}

func (d *TargetDispatcher) Scan(target models.Target) error {
	if target.Available == nil || !*target.Available || target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	return handleTarget(d, target)
}

func (d *TargetDispatcher) Clean(target models.Target) error {
	// TODO: to be implemented
	if !target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	return handleTarget(d, target)
}

func handleTarget(dispatcher *TargetDispatcher, target models.Target) error {

	item, _ := dispatcher.Bucket.SearchOrStore(&TargetDispatcherItem{
		id: target.ID,
	})

	item.Lock()
	defer item.Unlock()
	item.pending = &target

	if item.processing != nil {
		return nil
	}

	go func() {

		for {
			item.Lock()
			if item.pending != nil {
				item.processing = item.pending
			}
			if item.processing == nil {
				defer item.Unlock()
				break
			}
			target_ := item.processing
			item.Unlock()

			if target_.DeletedAt.Valid {
				_ = dispatcher.TargetService.Clean(target.ID)
			} else {
				_ = dispatcher.TargetService.Scan(target.ID)
			}

			// TODO: logging err
		}

	}()

	return nil
}
