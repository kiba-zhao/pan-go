package impl

import (
	"cmp"
	"pan/cache"
	"pan/extfs/dispatchers"
	"pan/extfs/errors"
	"pan/extfs/models"
	"pan/extfs/services"
	"sync"
)

type TargetDispatcherItem struct {
	id          uint
	pending     *models.Target
	processing  *models.Target
	done        dispatchers.DispatchDone
	pendingDone dispatchers.DispatchDone
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

func (d *TargetDispatcher) Scan(target models.Target, done dispatchers.DispatchDone) error {
	if target.Available == nil || !*target.Available || target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	return handleTarget(d, target, done)
}

func (d *TargetDispatcher) Clean(target models.Target, done dispatchers.DispatchDone) error {
	// TODO: to be implemented
	if !target.DeletedAt.Valid {
		return errors.ErrConflict
	}

	return handleTarget(d, target, done)
}

func handleTarget(dispatcher *TargetDispatcher, target models.Target, done dispatchers.DispatchDone) error {

	item, _ := dispatcher.Bucket.SearchOrStore(&TargetDispatcherItem{
		id: target.ID,
	})

	item.Lock()
	defer item.Unlock()
	item.pending = &target
	item.pendingDone = done

	if item.processing != nil {
		return nil
	}

	go func() {

		for {
			item.Lock()
			item.processing = item.pending
			item.done = item.pendingDone
			item.pending = nil
			item.pendingDone = nil
			if item.processing == nil {
				defer item.Unlock()
				break
			}
			target_ := item.processing
			item.Unlock()

			var err error
			if target_.DeletedAt.Valid {
				err = dispatcher.TargetService.Clean(target.ID)
			} else {
				err = dispatcher.TargetService.Scan(target.ID)
			}

			if item.done != nil {
				item.done(err)
			}
		}

	}()

	return nil
}
