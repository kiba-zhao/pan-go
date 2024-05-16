package impl

import (
	"cmp"
	"pan/app/cache"
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

type TargetDispatcherBucket = cache.Bucket[uint, *TargetDispatcherItem]

func NewTargetDispatcherBucket() TargetDispatcherBucket {
	bucket := cache.NewBucket[uint, *TargetDispatcherItem](cmp.Compare[uint])
	return cache.WrapSyncBucket(bucket)
}

type TargetDispatcher struct {
	TargetService *services.TargetService
	Bucket        TargetDispatcherBucket
}

func (d *TargetDispatcher) Scan(target models.Target) error {
	if target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	if !target.Available {
		return nil
	}

	return handleTarget(d, target)
}

func (d *TargetDispatcher) Clean(target models.Target) error {

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
			item.processing = item.pending
			item.pending = nil
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

		}

	}()

	return nil
}
