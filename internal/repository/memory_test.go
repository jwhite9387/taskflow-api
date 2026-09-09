package repository

import (
	"errors"
	"testing"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/service"
)

func TestMemoryUserRepository_Create(t *testing.T) {
	id := uuid.New()
	testUser := service.User{
		ID:           id,
		Username:     "test user",
		Email:        "test@test.com",
		PasswordHash: "test1234",
	}

	repo := NewMemoryUserRepository()
	err := repo.Create(testUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	storedUser, exists := repo.users[id]
	if !exists {
		t.Errorf("expected user to exist in repository")
	}
	if storedUser != testUser {
		t.Errorf("expected %v, got %v", testUser, storedUser)
	}
}

func TestMemoryUserRepository_CreateDuplicate(t *testing.T) {
	repo := NewMemoryUserRepository()
	id := uuid.New()
	testUser := service.User{
		ID:           id,
		Username:     "test user",
		Email:        "test@test.com",
		PasswordHash: "test1234",
	}

	err := repo.Create(testUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = repo.Create(testUser)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}
