package task

type inMemoryRepository struct {
	tasks map[int]*Task
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		tasks: make(map[int]*Task),
	}
}

func (r *inMemoryRepository) Create(externalID, cost, reward int) (*Task, error) {
	task := &Task{
		ID:          len(r.tasks) + 1,
		ExternalID:  externalID,
		Cost:        cost,
		Reward:      reward,
	}
	r.tasks[task.ID] = task
	return task, nil
}