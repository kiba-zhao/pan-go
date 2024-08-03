package impl

import (
	"cmp"
	"slices"

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

type TargetDispatcher struct {
	TargetService *services.TargetService
	items         []*TargetDispatcherItem
	locker        sync.Mutex
}

func NewTargetDispatcher() *TargetDispatcher {
	dispatcher := &TargetDispatcher{}
	dispatcher.items = make([]*TargetDispatcherItem, 0)
	return dispatcher
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

func compareTargetDispatcherItem(item *TargetDispatcherItem, key uint) int {
	return cmp.Compare(item.id, key)
}

func handleTarget(dispatcher *TargetDispatcher, target models.Target) error {
	dispatcher.locker.Lock()
	idx, ok := slices.BinarySearchFunc(dispatcher.items, target.ID, compareTargetDispatcherItem)
	var item *TargetDispatcherItem
	if ok {
		item = dispatcher.items[idx]
	} else {
		item = &TargetDispatcherItem{
			id: target.ID,
		}
		dispatcher.items = slices.Insert(dispatcher.items, idx, item)
	}
	dispatcher.locker.Unlock()

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
