// Define node search file task worker
//
// It is responsible for generating search files for node search tasks.
package nodesearchfile

import (
	"context"
	"errors"
	"os"
	nodeitem "pan/features/extfs/node_item"
	"pan/lib/log"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var ErrNodeSearchTaskAborted = errors.New("nodesearchfile.Worker Error: Search Task Aborted")
var ErrNoFileRater = errors.New("nodesearchfile.Worker Error: No File Rater")
var ErrGernerateRateFailed = errors.New("nodesearchfile.Worker Error: Generate Rate Failed")

type TaskWorker interface {
	RegisterFileRater(fileRater FileRater)
	UnregisterFileRater(fileRater FileRater)
	FileRaters() []FileRater
	Reload()
}

type stdTaskWorker struct {
	NodeItemService       nodeitem.NodeItemInternalService
	NodeSearchTaskRepo    NodeSearchTaskRepository
	NodeSearchFileService *NodeSearchFileService

	logger     log.Logger
	reloadChan chan struct{}
	reloadLock sync.Mutex
	reload     bool

	parallelThreshold   uint16
	parallelThresholdRW sync.RWMutex

	fileRaters   []FileRater
	fileRatersRW sync.RWMutex

	taskIds []uint64
	taskRW  sync.RWMutex
	taskNum atomic.Int32
}

func (worker *stdTaskWorker) ParallelThreshold() uint16 {
	worker.parallelThresholdRW.RLock()
	defer worker.parallelThresholdRW.RUnlock()
	return worker.parallelThreshold
}

func (worker *stdTaskWorker) SetParallelThreshold(threshold uint16) {
	worker.parallelThresholdRW.Lock()
	defer worker.parallelThresholdRW.Unlock()
	worker.parallelThreshold = threshold
	worker.Reload()
}

func (worker *stdTaskWorker) FileRaters() []FileRater {
	worker.fileRatersRW.RLock()
	defer worker.fileRatersRW.RUnlock()
	return worker.fileRaters
}

func (worker *stdTaskWorker) RegisterFileRater(fileRater FileRater) {
	worker.fileRatersRW.Lock()
	defer worker.fileRatersRW.Unlock()

	worker.fileRaters = append(worker.fileRaters, fileRater)
}

func (worker *stdTaskWorker) UnregisterFileRater(fileRater FileRater) {
	worker.fileRatersRW.Lock()
	defer worker.fileRatersRW.Unlock()

	for idx, rater := range worker.fileRaters {
		if rater == fileRater {
			worker.fileRaters = slices.Delete(worker.fileRaters, idx, idx+1)
			break
		}
	}
}

func (worker *stdTaskWorker) HasTask(taskId uint64) bool {
	worker.taskRW.RLock()
	defer worker.taskRW.RUnlock()
	_, ok := slices.BinarySearch(worker.taskIds, taskId)
	return ok
}

func (worker *stdTaskWorker) PickTask() (*NodeSearchTask, error) {
	worker.taskRW.Lock()
	defer worker.taskRW.Unlock()

	taskIds := worker.taskIds
	task, err := worker.NodeSearchTaskRepo.SelectWithStatusExcludeIds(NodeSearchTaskStatusPending, taskIds)
	if err != nil {
		return nil, err
	}

	if taskIds == nil {
		taskIds = make([]uint64, 0)
		taskIds = append(taskIds, task.ID)
	} else {
		idx, _ := slices.BinarySearch(taskIds, task.ID)
		taskIds = slices.Insert(taskIds, idx, task.ID)
	}
	worker.taskIds = taskIds

	return &task, nil
}

func (worker *stdTaskWorker) ReleaseTask(taskId uint64) {
	worker.taskRW.Lock()
	defer worker.taskRW.Unlock()
	taskIds := worker.taskIds
	if len(taskIds) <= 0 {
		return
	}
	if idx, ok := slices.BinarySearch(taskIds, taskId); ok {
		worker.taskIds = slices.Delete(taskIds, idx, idx+1)
	}
}

func (worker *stdTaskWorker) Reload() {
	worker.logger.Debug("nodesearchfile.TaskWorker", "Reload")

	worker.reloadLock.Lock()
	defer worker.reloadLock.Unlock()
	if worker.reload {
		return
	}

	worker.reload = true
	worker.reloadChan <- struct{}{}
}

func (worker *stdTaskWorker) Run(ctx context.Context) error {
	worker.logger.Debug("nodesearchfile.TaskWorker", "Run begin")
	defer worker.logger.Debug("nodesearchfile.TaskWorker", "Run end")

	var err error
	var closed bool

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-worker.reloadChan:
			worker.reloadLock.Lock()
			worker.reload = false
			worker.reloadLock.Unlock()
		}

		if closed {
			break
		}

		parallelThreshold := int32(worker.ParallelThreshold())
		for taskNum := worker.taskNum.Load(); taskNum < parallelThreshold; taskNum++ {
			go worker.RunTask(ctx)
		}

	}
	return err
}

func (worker *stdTaskWorker) RunTask(ctx context.Context) error {

	var err error
	var closed bool
	var isAlive bool
	timeCh := time.After(0)

	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			closed = true
		case <-timeCh:
		}

		if closed {
			break
		}

		parallelThreshold := int32(worker.ParallelThreshold())
		taskNum := worker.taskNum.Load()
		if isAlive {
			if taskNum > parallelThreshold {
				if ok := worker.taskNum.CompareAndSwap(taskNum, taskNum-1); ok {
					break
				}
				timeCh = time.After(1 * time.Second)
				continue
			}
		} else {
			if taskNum >= parallelThreshold {
				break
			}
			if ok := worker.taskNum.CompareAndSwap(taskNum, taskNum+1); !ok {
				timeCh = time.After(1 * time.Second)
				continue
			}
			isAlive = true
		}

		task, err := worker.PickTask()
		if err != nil {
			if errors.Is(err, ErrNodeSearchTaskNotFound) {
				worker.taskNum.Add(-1)
				break
			}
			timeCh = time.After(5 * time.Second)
			continue
		}

		for failed := 0; ; failed++ {
			err = worker.RushTask(ctx, task)
			if err == nil || err == ctx.Err() {
				err = nil
				break
			}
			worker.logger.Error("TaskWorker", "RunTask Error"+err.Error())
			if failed >= 5 {
				break
			}
			select {
			case <-ctx.Done():
				break
			case <-time.After(1 * time.Second):
			}
		}

		if err != nil {
			task.Status = NodeSearchTaskStatusError
			worker.NodeSearchTaskRepo.UpdateWithStatus(NodeSearchTaskStatusPending, *task)
		}

		worker.ReleaseTask(task.ID)
		timeCh = time.After(0)
	}

	return err

}

func (worker *stdTaskWorker) RushTask(ctx context.Context, task *NodeSearchTask) error {

	items, err := worker.NodeItemService.SelectAllWithEnabled(true)
	if err != nil {
		return err
	}

	raters := worker.FileRaters()
	if len(raters) <= 0 {
		return ErrNoFileRater
	}

	err = worker.NodeSearchFileService.InitWithTaskID(task.ID)
	if err != nil {
		return err
	}

	ratersTokens := make([]Tokens, 0)
	for _, rater := range raters {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tokens := rater.Tokenize(task.Query)
		ratersTokens = append(ratersTokens, tokens)
	}

	finishedCount := 0

	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		root := item.FilePath
		filePaths, err := WalkRoot(root)
		if err != nil {
			continue
		}

		for filePath := range filePaths {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			rate, err := worker.NewNodeSearchFile(item, filePath, raters, ratersTokens)
			if err == nil {
				_, err = worker.NodeSearchFileService.SaveWithTaskID(task.ID, rate)
			}

		}
		finishedCount++
	}

	if finishedCount == len(items) {
		task.Status = NodeSearchTaskStatusSuccess
	} else {
		task.Status = NodeSearchTaskStatusWarning
	}

	_, err = worker.NodeSearchTaskRepo.UpdateWithStatus(NodeSearchTaskStatusPending, *task)

	return err
}

func (w *stdTaskWorker) NewNodeSearchFile(item nodeitem.NodeItem, filePath string, raters []FileRater, ratersTokens []Tokens) (NodeSearchFile, error) {

	var rate NodeSearchFile
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
