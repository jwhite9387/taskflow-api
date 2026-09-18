package repository

import (
	"errors"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/service"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type MemoryUserRepository struct {
	users map[uuid.UUID]service.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[uuid.UUID]service.User),
	}
}

func (r *MemoryUserRepository) Create(user service.User) error {
	_, exists := r.users[user.ID]
	if exists {
		return ErrUserAlreadyExists
	}

	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) FindByEmail(email string) (service.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return service.User{}, service.ErrUserNotFound
}

func (r *MemoryUserRepository) FindByID(id uuid.UUID) (service.User, error) {
	user, ok := r.users[id]

	if ok {
		return user, nil
	}
	return service.User{}, service.ErrUserNotFound
}
