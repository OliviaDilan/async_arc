package task

type Task struct {
	ID          int
	Cost        int
	Reward      int
}

type Repository interface {
	Create(externalID, cost, reward int) (*Task, error)
	GetTaskByExternalID(externalID int) (*Task, error)
}