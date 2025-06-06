// Define node search file task cleaner
package nodesearchfile

import (
	"context"
	"pan/lib/log"
	"sync/atomic"
	"time"
)

type stdTaskCleaner struct {
	NodeSearchTaskRepo    NodeSearchTaskRepository
	NodeSearchFileService *NodeSearchFileService

	logger    log.Logger
	worker    *stdTaskWorker
	lifecycle atomic.Uint64
}

func (cleaner *stdTaskCleaner) SetLifecycle(lifecycle uint64) {
	cleaner.lifecycle.Swap(lifecycle)
}

func (cleaner *stdTaskCleaner) Lifecycle() uint64 {
	return cleaner.lifecycle.Load()
}

func (cleaner *stdTaskCleaner) Run(ctx context.Context) error {
	cleaner.logger.Debug("nodesearchfile.TaskCleaner", "Run begin")
	defer cleaner.logger.Debug("nodesearchfile.TaskCleaner", "Run end")

	var err error
	timeCh := time.After(0)

read_loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break read_loop
		case <-timeCh:
		}

		lifecycle := cleaner.Lifecycle()
		tasks, err := cleaner.NodeSearchTaskRepo.SearchWithLifecycle(lifecycle)
		if err != nil || len(tasks) <= 0 {
			timeCh = time.After(time.Duration(lifecycle))
			continue
		}

		hasDestroy := false
		for _, task := range tasks {
			if cleaner.worker.HasTask(task.ID) {
				continue
			}
			hasDestroy = true
			taskErr := cleaner.NodeSearchFileService.DestroyWithTaskID(task.ID)
			if taskErr != nil {
				continue
			}
			cleaner.NodeSearchTaskRepo.Delete(task)
		}

		if hasDestroy {
			timeCh = time.After(time.Duration(lifecycle))
		} else {
			timeCh = time.After(0)
		}
	}

	return err
}
