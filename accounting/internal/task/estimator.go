package task

type Task struct {
	ID          int
	ExternalID  int
	Cost        int
	Reward      int
}

type Repository interface {
	Create(externalID, cost, reward int) (*Task, error)
}