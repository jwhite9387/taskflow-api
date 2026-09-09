package service

import (
	"strings"
	"uuid"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
}

type UserRepository interface {
	Create(user User) error
}

type UserService struct {
	repo UserRepository
}

type ValidationError struct {
	Errors []string
}

func (e ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidateRegistration(username, email, password string) error {
	var errs []string
	if username == "" {
		errs = append(errs, "Username is required")
	}
	if email == "" {
		errs = append(errs, "Email is required")
	}
	if password == "" {
		errs = append(errs, "Password is required")
	} else if len(password) < 8 {
		errs = append(errs, "Password must be at least 8 characters")
	}

	if len(errs) > 0 {
		return ValidationError{
			Errors: errs,
		}
	}
	return nil
}

func (s *UserService) Register(username, email, password string) error {
	if err := ValidateRegistration(username, email, password); err != nil {
		return err
	}
	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}
	id := uuid.New()
	user := User{
		ID:           id,
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}
	if err := s.repo.Create(user); err != nil {
		return err
	}
	return nil
}
