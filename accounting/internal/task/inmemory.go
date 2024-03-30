package task

import (
	"fmt"
)

type inMemoryRepository struct {
	tasks map[int]*Task // map[taskID]task
	externalIDIndex map[int]int // map[externalID]taskID
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		tasks: make(map[int]*Task),
		externalIDIndex: make(map[int]int),
	}
}

func (r *inMemoryRepository) Create(externalID, cost, reward int) (*Task, error) {
	task := &Task{
		ID:          len(r.tasks) + 1,
		Cost:        cost,
		Reward:      reward,
	}
	r.tasks[task.ID] = task
	r.externalIDIndex[externalID] = task.ID
	return task, nil
}

func (r *inMemoryRepository) GetTaskByExternalID(externalID int) (*Task, error) {
	task, ok := r.tasks[r.externalIDIndex[externalID]]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}