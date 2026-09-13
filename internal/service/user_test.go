package service

import (
	"errors"
	"testing"
	"uuid"
)

type fakeUserRepository struct {
	user         User
	createCalled bool
	err          error
}

func (f *fakeUserRepository) Create(user User) error {
	f.createCalled = true
	f.user = user
	return f.err
}

func (f *fakeUserRepository) FindByEmail(email string) (User, error) {
	return User{}, ErrUserNotFound
}

func TestUserService_Register(t *testing.T) {
	repo := &fakeUserRepository{}
	service := NewUserService(repo)
	password := "password123"

	if err := service.Register("josh", "josh@example.com", password); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.user.ID == (uuid.UUID{}) {
		t.Fatalf("expected user ID to be generated")
	}

	if !VerifyPassword(password, repo.user.PasswordHash) {
		t.Fatalf("expected password to verify against stored hash")
	}

	if repo.user.PasswordHash == password {
		t.Fatalf("expected password to be hashed, got plaintext password")
	}
}

func TestUserService_RegisterValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		email    string
		password string
	}{
		{"short password", "josh", "josh@example.com", "short"},
		{"missing username", "", "josh@example.com", "password123"},
		{"missing email", "josh", "", "password123"},
		{"missing password", "josh", "josh@example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepository{}
			service := NewUserService(repo)
			if err := service.Register(tt.username, tt.email, tt.password); err == nil {
				t.Fatalf("expected error for invalid registration")
			}

			if repo.createCalled {
				t.Fatalf("expected repository Create not to be called")
			}
		})
	}
}

func TestUserService_RegisterRepositoryError(t *testing.T) {
	repoErr := errors.New("repository error")

	repo := &fakeUserRepository{
		err: repoErr,
	}

	service := NewUserService(repo)
	err := service.Register("josh", "josh@example.com", "password123")

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestUserService_LoginUserNotFound(t *testing.T) {
	repo := &fakeUserRepository{}
	testService := NewUserService(repo)

	_, err := testService.Login("doesnotexist@example.com", "password123")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
