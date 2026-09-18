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

type fakeTokenCreator struct{}

func (f *fakeUserRepository) Create(user User) error {
	f.createCalled = true
	f.user = user
	return f.err
}

func (f *fakeUserRepository) FindByEmail(email string) (User, error) {
	if f.user.Email != email {
		return User{}, ErrUserNotFound
	}
	return f.user, nil
}

func (f *fakeUserRepository) FindByID(id uuid.UUID) (User, error) {
	if f.user.ID != id {
		return User{}, ErrUserNotFound
	}
	return f.user, nil
}

func (f *fakeTokenCreator) CreateToken(userID uuid.UUID) (string, error) {
	return "test-token", nil
}

func TestUserService_Register(t *testing.T) {
	repo := &fakeUserRepository{}
	tokenCreator := &fakeTokenCreator{}
	testService := NewUserService(repo, tokenCreator)
	password := "password123"

	if err := testService.Register("josh", "josh@example.com", password); err != nil {
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
			tokenCreator := &fakeTokenCreator{}
			testService := NewUserService(repo, tokenCreator)
			if err := testService.Register(tt.username, tt.email, tt.password); err == nil {
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
	tokenCreator := &fakeTokenCreator{}

	testService := NewUserService(repo, tokenCreator)
	err := testService.Register("josh", "josh@example.com", "password123")

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestUserService_LoginUserNotFound(t *testing.T) {
	repo := &fakeUserRepository{}
	tokenCreator := &fakeTokenCreator{}
	testService := NewUserService(repo, tokenCreator)

	_, err := testService.Login("doesnotexist@example.com", "password123")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserService_LoginSuccess(t *testing.T) {
	password := "password123"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}
	repo := &fakeUserRepository{
		user: testUser,
	}
	tokenCreator := &fakeTokenCreator{}

	testService := NewUserService(repo, tokenCreator)

	result, err := testService.Login("test@example.com", password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Token != "test-token" {
		t.Fatalf("expected token %q, got %q", "test-token", result.Token)
	}
}

func TestUserService_LoginInvalidPassword(t *testing.T) {
	password := "password123"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}
	repo := &fakeUserRepository{
		user: testUser,
	}
	tokenCreator := &fakeTokenCreator{}

	testService := NewUserService(repo, tokenCreator)

	_, err = testService.Login("test@example.com", "wrongpassword")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}
	repo := &fakeUserRepository{
		user: testUser,
	}

	tokenCreator := &fakeTokenCreator{}
	testService := NewUserService(repo, tokenCreator)

	result, err := testService.GetUserByID(testUser.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != testUser {
		t.Fatalf("expected user %v, got %v", testUser, result)
	}
}

func TestUserService_GetUserByIDNotFound(t *testing.T) {
	repo := &fakeUserRepository{}
	testService := NewUserService(repo, &fakeTokenCreator{})

	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = testService.GetUserByID(userID)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
