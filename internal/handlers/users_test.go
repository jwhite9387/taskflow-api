package handlers

import (
	"bytes"
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

type fakeTokenCreator struct{}

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

func (f *fakeTokenCreator) CreateToken(userID uuid.UUID) (string, error) {
	return "test-token", nil
}

func TestUserHandler_InvalidPassword(t *testing.T) {
	repo := &fakeUserRepository{}
	tokenCreator := &fakeTokenCreator{}
	testService := service.NewUserService(repo, tokenCreator)
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

func TestUserHandler_UserNotFound(t *testing.T) {
	repo := &fakeUserRepository{}
	tokenCreator := &fakeTokenCreator{}
	testService := service.NewUserService(repo, tokenCreator)

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

func TestUserHandler_LoginSuccess(t *testing.T) {
	password := "password123"
	passwordHash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testUser := service.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}

	repo := &fakeUserRepository{
		user: testUser,
	}

	tokenCreator := &fakeTokenCreator{}

	testService := service.NewUserService(repo, tokenCreator)
	handler := &UserHandler{
		userService: testService,
	}

	loginRequest := LoginUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	buffer := bytes.NewBuffer(body)

	request := httptest.NewRequest(
		http.MethodPost,
		"/users/login",
		buffer,
	)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handler.LoginUser(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response LoginUserResponse

	err = json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("expected valid JSON response, got %v", err)
	}

	if response.Token != "test-token" {
		t.Fatalf("expected token %q, got %q", "test-token", response.Token)
	}

	if response.ID != userID {
		t.Fatalf("expected ID %q, got %q", userID, response.ID)
	}

	if response.Username != "testuser" {
		t.Fatalf("expected username %q, got %q", "testuser", response.Username)
	}

	if response.Email != "test@example.com" {
		t.Fatalf("expected email %q, got %q", "test@example.com", response.Email)
	}

}
