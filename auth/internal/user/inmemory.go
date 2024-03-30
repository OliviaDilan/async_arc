package user

import (
	"fmt"
)

type inMemoryRepository struct {
	users map[string]*User
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *inMemoryRepository) Create(user *User) error {
	r.users[user.Username] = user
	return nil
}

func (r *inMemoryRepository) Get(username string, password string) (*User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	if user.Password != password {
		return nil, fmt.Errorf("wrong password")
	}
	return user, nil
}

func (r *inMemoryRepository) GetAll() ([]*User, error) {
	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}