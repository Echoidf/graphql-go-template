package task

import "sync"

type TaskManager struct {
	tasks map[string]*DelayedTask
	lock  sync.Mutex
}

func NewTaskManager() *TaskManager {
	return &TaskManager{}
}

func (tm *TaskManager) AddTask(id string, task *DelayedTask) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if tm.tasks == nil {
		tm.tasks = make(map[string]*DelayedTask)
	}
	tm.tasks[id] = task
}

func (tm *TaskManager) CancelAll() {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	for _, task := range tm.tasks {
		task.Cancel()
	}
	tm.tasks = nil
}

func (tm *TaskManager) CancelTask(id string) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if task, ok := tm.tasks[id]; ok {
		task.Cancel()
		delete(tm.tasks, id)
	}
}
