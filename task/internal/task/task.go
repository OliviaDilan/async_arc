package task

type Task struct {
	ID          int
	Title       string
	Assignee    string
	Status      Status
}

const (
	StatusOpen   Status = "Open"
	StatusClosed Status = "Closed"
)

type Status string

type Repository interface {
	Create(name string) (*Task, error)
	GetAll() ([]*Task, error)
	GetByAssignee(assignee string) ([]*Task, error)
	GetByStatus(status Status) ([]*Task, error)
	GetByID(taskID int) (*Task, error)
	Close(taskID int) error
	Assign(taskID int, username string) error
}