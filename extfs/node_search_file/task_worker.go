// Define node search file task worker
//
// It is responsible for generating search files for node search tasks.
package nodesearchfile

import (
	"context"
	"errors"
	"os"
	"pan/app/web"
	nodeitem "pan/extfs/node_item"
	"strings"
	"sync"
	"time"
)

var ErrNodeSearchTaskAborted = errors.New("nodesearchfile.Worker Error: Search Task Aborted")
var ErrNoFileRater = errors.New("nodesearchfile.Worker Error: No File Rater")
var ErrGernerateRateFailed = errors.New("nodesearchfile.Worker Error: Generate Rate Failed")

type TaskWorker interface {
	// Reload triggers the reloading process for the task worker. It ensures that
	// the necessary components are reinitialized or refreshed to reflect any updates
	// or changes. This method is typically used when configurations or dependencies
	// are updated and the worker needs to adapt to those changes.
	Reload()
}

type taskWorkerImpl struct {
	NodeItemService       nodeitem.NodeItemInternalService
	NodeSearchTaskRepo    NodeSearchTaskRepository
	NodeSearchFileService *NodeSearchFileService
	Agent                 Agent

	reloadLocker sync.RWMutex
	reload       bool
	reloadChan   chan struct{}
	reloadOnce   sync.Once

	tasks       []NodeSearchTask
	tasksLocker sync.RWMutex
}

// ReloadChan returns a channel that receives a signal when the task worker needs to be reloaded.
//
// It is used by the task scheduler to reload the task worker.
func (w *taskWorkerImpl) ReloadChan() chan struct{} {
	w.reloadOnce.Do(func() {
		w.reloadChan = make(chan struct{}, 1)
	})
	return w.reloadChan
}

// Reload triggers the reloading process for the task worker.
//
// It locks the reloadLocker to ensure thread-safety, checks if the reload flag is already set,
// and if not, sets the flag to true and sends a signal to the reload channel to initiate the reload.

func (w *taskWorkerImpl) Reload() {
	w.reloadLocker.Lock()
	defer w.reloadLocker.Unlock()
	if w.reload {
		return
	}
	w.reload = true
	w.ReloadChan() <- struct{}{}
}

func (w *taskWorkerImpl) Ready(ctx context.Context) error {

	var tasks []NodeSearchTask
	var err error

	total := int64(0)
	closed := false

	for {

		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-w.ReloadChan():
		}

		w.reloadLocker.Lock()
		if !closed {
			settings := w.Agent.Settings()

			var condition web.RangeCondition
			condition.RangeEnd = int(settings.TaskNum)
			condition.RangeStart = 0

			total, tasks, err = w.NodeSearchTaskRepo.SearchWithStatus(NodeSearchTaskStatusPending, condition)
			if err != nil {
				w.reloadLocker.Unlock()
				<-time.After(5 * time.Second)
				continue
			}

			w.reload = total > 0 && total <= int64(len(tasks))
			if w.reload {
				w.reloadChan <- struct{}{}
			}
		}
		w.reloadLocker.Unlock()

		if closed {
			break
		}

		if len(tasks) <= 0 {
			continue
		}

		err = w.RunTasks(tasks)
		if err != nil {
			w.Reload()
		}

	}
	return err
}

// RunTasks runs the given node search tasks. It retrieves node items, uses file raters to
// generate ratings for the items, and saves the results to the database. It returns an error
// if any issues occur during the process. If a task is aborted, it returns ErrNodeSearchTaskAborted.
// If no file rater is available, it returns ErrNoFileRater. If generating a rating fails, it
// returns ErrGernerateRateFailed.
func (w *taskWorkerImpl) RunTasks(tasks []NodeSearchTask) error {

	var err error
	w.tasksLocker.Lock()
	w.tasks = tasks
	w.tasksLocker.Unlock()

	for idx, task := range tasks {
		items, err := w.NodeItemService.SelectAllWithEnabled(true)
		if err != nil {
			return err
		}

		raters := w.Agent.FileRaters()
		if len(raters) <= 0 {
			return ErrNoFileRater
		}

		err = w.NodeSearchFileService.InitWithTaskID(task.ID)
		if err != nil {
			return err
		}

		ratersTokens := make([]Tokens, 0)
		for _, rater := range raters {
			tokens := rater.Tokenize(task.Query)
			ratersTokens = append(ratersTokens, tokens)
		}

		finishedCount := 0
	items_loop:
		for _, item := range items {
			root := item.FilePath
			filePaths, err := WalkRoot(root)
			if err != nil {
				continue
			}

			for filePath := range filePaths {
				rate, err := w.NewNodeSearchFile(idx, item, filePath, raters, ratersTokens)
				if errors.Is(err, ErrNodeSearchTaskAborted) {
					break items_loop
				}
				if err != nil {
					continue
				}

				_, err = w.NodeSearchFileService.SaveWithTaskID(task.ID, rate)
			}
			finishedCount++
		}

		if finishedCount < 0 {
			continue
		}

		if finishedCount == len(items) {
			task.Status = NodeSearchTaskStatusSuccess
		} else {
			task.Status = NodeSearchTaskStatusWarning
		}

		_, err = w.NodeSearchTaskRepo.UpdateWithStatus(NodeSearchTaskStatusPending, task)
	}

	return err
}

// NewNodeSearchFile creates a new NodeSearchFile object based on the given item and
// file path. It also rates the file using the given raters and tokens.
//
// It returns the new NodeSearchFile object and an error if any issues occur. If the
// task is not pending, it returns ErrNodeSearchTaskAborted. If no file rater is
// available, it returns ErrNoFileRater. If the file rating fails, it returns
// ErrGernerateRateFailed.
func (w *taskWorkerImpl) NewNodeSearchFile(taskIdx int, item nodeitem.NodeItem, filePath string, raters []FileRater, ratersTokens []Tokens) (NodeSearchFile, error) {

	var rate NodeSearchFile

	w.tasksLocker.RLock()
	task := w.tasks[taskIdx]
	if task.Status != NodeSearchTaskStatusPending {
		w.tasksLocker.RUnlock()
		return rate, ErrNodeSearchTaskAborted
	}
	w.tasksLocker.RUnlock()

	stat, err := os.Stat(filePath)
	if err != nil {
		return rate, err
	}

	if strings.Compare(item.FilePath, filePath) != 0 {
		rate.Name = stat.Name()
		rate.FilePath = filePath[len(item.FilePath):]
	} else {
		rate.Name = item.Name
	}

	rate.ItemID = item.ID
	rate.Size = stat.Size()
	rate.UpdatedAt = stat.ModTime()
	rate.CreatedAt = stat.ModTime()

	if stat.IsDir() {
		rate.FileType = nodeitem.FileTypeFolder
	} else {
		rate.FileType = nodeitem.FileTypeFile
	}

	for idx, rater := range raters {
		tokens := ratersTokens[idx]
		score, err := rater.Rate(filePath, tokens)
		if err != nil {
			continue
		}
		if score > 0 && score > rate.Score {
			rate.Score = score
			rate.Tokens = tokens
		}
	}

	if rate.Score <= 0 {
		return rate, ErrGernerateRateFailed
	}

	mimeType, err := nodeitem.GenerateMimeTypeWithFilePath(filePath)
	if err == nil {
		rate.MimeType = mimeType
	}

	return rate, nil
}
