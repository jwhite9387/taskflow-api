package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/auth"
	"github.com/jwhite9387/taskflow-api/internal/service"
)

func TestUsersMeRoute_Unauthorized(t *testing.T) {
	repo := &fakeUserRepository{}
	tokenCreator := &fakeTokenCreator{}
	userService := service.NewUserService(repo, tokenCreator)

	jwtManager := auth.NewJWTManager("test-secret")

	mux := http.NewServeMux()
	RegisterRoutes(mux, userService, jwtManager)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}

func TestUsersMeRoute_Authenticated(t *testing.T) {
	userID := uuid.New()

	testUser := service.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	repo := &fakeUserRepository{
		user: testUser,
	}

	tokenCreator := &fakeTokenCreator{}
	userService := service.NewUserService(repo, tokenCreator)

	jwtManager := auth.NewJWTManager("test-secret")

	token, err := jwtManager.CreateToken(userID)
	if err != nil {
		t.Fatalf("failed to create token: %v", err)
	}

	mux := http.NewServeMux()
	RegisterRoutes(mux, userService, jwtManager)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}

	var response CurrentUserResponse

	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.ID != userID {
		t.Errorf("expected ID %v, got %v", userID, response.ID)
	}

	if response.Username != "testuser" {
		t.Errorf("expected username testuser, got %s", response.Username)
	}

	if response.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", response.Email)
	}
}
