package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/service"
)

type fakeUserRepository struct {
	user service.User
}

func (f *fakeUserRepository) FindByEmail(email string) (service.User, error) {
	if email == f.user.Email {
		return f.user, nil
	}

	return service.User{}, service.ErrUserNotFound
}

func (f *fakeUserRepository) Create(user service.User) error {
	f.user = user
	return nil
}

func TestUserHandler_LoginUser(t *testing.T) {
	repo := &fakeUserRepository{}
	testService := service.NewUserService(repo)
	password := "password123"

	passwordHash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := service.User{
		ID:           uuid.New(),
		Username:     "josh",
		Email:        "josh@example.com",
		PasswordHash: passwordHash,
	}
	repo.user = testUser

	handler := &UserHandler{
		userService: testService,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		strings.NewReader(`{"email":"josh@example.com","password":"password123"}`),
	)

	recorder := httptest.NewRecorder()

	handler.LoginUser(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response LoginUserResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("expected valid JSON response, got %v", err)
	}

	if response.ID != testUser.ID {
		t.Errorf("expected ID %v, got %v", testUser.ID, response.ID)
	}

	if response.Username != testUser.Username {
		t.Errorf("expected username %q, got %q", testUser.Username, response.Username)
	}

	if response.Email != testUser.Email {
		t.Errorf("expected email %q, got %q", testUser.Email, response.Email)
	}
}

func TestUserHandler_InvalidPassword(t *testing.T) {
	repo := &fakeUserRepository{}
	testService := service.NewUserService(repo)
	password := "password123"

	passwordHash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := service.User{
		ID:           uuid.New(),
		Username:     "josh",
		Email:        "josh@example.com",
		PasswordHash: passwordHash,
	}
	repo.user = testUser

	handler := &UserHandler{
		userService: testService,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		strings.NewReader(`{"email":"josh@example.com","password":"wrongpassword"}`),
	)

	recorder := httptest.NewRecorder()

	handler.LoginUser(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "invalid credentials") {
		t.Fatalf("expected invalid credentials message, got %q", recorder.Body.String())
	}
}

func TestHandler_UserNotFound(t *testing.T) {
	repo := &fakeUserRepository{}
	testService := service.NewUserService(repo)

	handler := &UserHandler{
		userService: testService,
	}

	request := httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		strings.NewReader(`{"email":"doesnotexist@example.com","password":"password123"}`),
	)

	recorder := httptest.NewRecorder()
	handler.LoginUser(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
