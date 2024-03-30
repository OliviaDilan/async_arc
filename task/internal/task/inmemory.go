package task

import (
	"fmt"
)

type inMemoryRepository struct {
	tasks map[int]*Task
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		tasks: make(map[int]*Task),
	}
}

func (r *inMemoryRepository) Create(name string) (*Task, error) {
	task := &Task{
		ID:     len(r.tasks) + 1,
		Title:  name,
		Status: StatusOpen,
	}
	r.tasks[task.ID] = task
	return task, nil
}

func (r *inMemoryRepository) GetAll() ([]*Task, error) {
	tasks := make([]*Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *inMemoryRepository) GetByAssignee(assignee string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	for _, task := range r.tasks {
		if task.Assignee == assignee {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *inMemoryRepository) GetByStatus(status Status) ([]*Task, error) {
	tasks := make([]*Task, 0)
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *inMemoryRepository) Close(taskID int) error {
	task, ok := r.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found")
	}
	task.Status = StatusClosed
	return nil
}

func (r *inMemoryRepository) Assign(taskID int, username string) error {
	task, ok := r.tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found")
	}
	task.Assignee = username
	return nil
}