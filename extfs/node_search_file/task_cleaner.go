// Define node search file task cleaner
package nodesearchfile

import (
	"context"
	"time"
)

type taskCleanerImpl struct {
	NodeSearchTaskRepo    NodeSearchTaskRepository
	NodeSearchFileService *NodeSearchFileService
	Agent                 Agent
}

// Ready executes the cleanup routine for node search file tasks based on their
// lifecycle. It continuously checks for tasks that have exceeded their lifecycle
// duration and deletes their associated search files. The function listens for
// context cancellation to gracefully terminate the cleanup process. It returns
// an error if the context is canceled or if an error occurs during task
// retrieval or deletion.

func (c *taskCleanerImpl) Ready(ctx context.Context) error {
	var err error
	var timeCh <-chan time.Time

read_loop:
	for {
		if timeCh == nil {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				break read_loop
			default:
			}
		} else {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				break read_loop
			case <-timeCh:
			}
			timeCh = nil
		}

		settings := c.Agent.Settings()
		tasks, err := c.NodeSearchTaskRepo.SearchWithLifecycle(settings.Lifecycle)
		if err != nil || len(tasks) <= 0 {
			timeCh = time.After(time.Duration(settings.Lifecycle))
			continue
		}

		for _, task := range tasks {
			taskErr := c.NodeSearchFileService.DestroyWithTaskID(task.ID)
			if taskErr != nil {
				continue
			}
			c.NodeSearchTaskRepo.Delete(task)
		}
	}

	return err
}
