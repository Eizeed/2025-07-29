package task

import (
	"errors"
	"sync"

	"github.com/Eizeed/2025-07-29/internal/pkg/archive"
	"github.com/Eizeed/2025-07-29/internal/pkg/constants"
	"github.com/Eizeed/2025-07-29/pkg/uuid"
)

const (
	ErrQueueFull = "ErrQueueFull"
	ErrTaskFull  = "ErrTaskFull"
)

type TaskQueue struct {
	inner []Task
	rw    sync.RWMutex
}

func NewQueue() TaskQueue {
	return TaskQueue{
		inner: []Task{},
		rw:    sync.RWMutex{},
	}
}

func (queue *TaskQueue) ViewTasks() []Task {
	return queue.inner
}

func (queue *TaskQueue) InsertTask() (uuid.UUID, error) {
	queue.rw.Lock()
	defer queue.rw.Unlock()

	if len(queue.inner) >= constants.TASK_LIMIT {
		return uuid.UUID{}, errors.New(ErrQueueFull)
	}

	task := NewTask()

	queue.inner = append(queue.inner, task)

	return task.UUID, nil
}

func (queue *TaskQueue) RemoveByUUID(uuid uuid.UUID) (Task, bool) {
	queue.rw.Lock()
	defer queue.rw.Unlock()

	for i, task := range queue.inner {
		if task.UUID == uuid {
			queue.inner = append(queue.inner[:i], queue.inner[i+1:]...)
			return task, true
		}
	}

	return Task{}, false
}

func (queue *TaskQueue) GetTask(uuid uuid.UUID) (*Task, bool) {
	queue.rw.RLock()
	defer queue.rw.RUnlock()

	for i, task := range queue.inner {
		if task.UUID == uuid {
			return &queue.inner[i], true
		}
	}

	return nil, false
}

type Task struct {
	UUID    uuid.UUID
	Archive archive.Archive
	mutex   *sync.Mutex
}

func NewTask() Task {
	uuid := uuid.NewV4()

	return Task{
		UUID:    uuid,
		Archive: archive.NewArchive(),
		mutex:   &sync.Mutex{},
	}
}

func (task *Task) Push(path string) error {
	task.mutex.Lock()
	defer task.mutex.Unlock()

	err := task.Archive.AddPath(path)
	if err != nil {
		return err
	}

	return nil
}
