package user

type User struct {
	Username string
	Password string
	Role     string
}

type Repository interface {
	Create(user *User) error
	Get(username string, password string) (*User, error)
	GetAll() ([]*User, error)
}