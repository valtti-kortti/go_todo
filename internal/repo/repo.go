package repo

import (
	"sync"
)

type repository struct {
	mu    sync.RWMutex
	tasks map[int]*Task
	id    int
}

type Repository interface {
	CreateNewTask(task *Task) int
	GetAllTasks() []*Task
	GetTaskById(taskId int) (*Task, bool)
	DeleteTaskById(taskId int) bool
	UpdateTaskById(taskId int, data Task) *Task
}

func NewRepository() Repository {
	return &repository{tasks: make(map[int]*Task), id: 0}
}

func (r *repository) CreateNewTask(task *Task) int {
	task.Status = "New"
	r.mu.Lock()
	defer r.mu.Unlock()
	r.id += 1
	r.tasks[r.id] = task
	return r.id
}

func (r *repository) GetAllTasks() []*Task {
	result := make([]*Task, 0, len(r.tasks))
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, v := range r.tasks {
		result = append(result, v)
	}
	return result
}

func (r *repository) GetTaskById(taskId int) (*Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.tasks[taskId]; !ok {
		return nil, false
	}
	return r.tasks[taskId], true
}

func (r *repository) DeleteTaskById(taskId int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[taskId]; !ok {
		return false
	}
	delete(r.tasks, taskId)
	return true
}

func (r *repository) UpdateTaskById(taskId int, input Task) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.tasks[taskId]
	if !ok {
		return nil
	}

	if input.Title != "" {
		task.Title = input.Title
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if input.Status != "" {
		task.Status = input.Status
	}

	return task
}
