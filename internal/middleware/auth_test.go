package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/auth"
)

func TestMiddlewareAuth_MissingHeader(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not have been called")
	})

	handler := RequireAuth(jwtManager, next)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %v", recorder.Code)
	}
}

func TestMiddlewareAuth_InvalidScheme(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not have been called")
	})

	handler := RequireAuth(jwtManager, next)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	request.Header.Set("Authorization", "Basic abc123")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %v", recorder.Code)
	}
}

func TestMiddlewareAuth_InvalidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret")

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("next handler should not have been called")
	})

	handler := RequireAuth(jwtManager, next)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	request.Header.Set("Authorization", "Bearer invalid-token")

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %v", recorder.Code)
	}
}

func TestMiddlewareAuth_ValidToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret")

	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	testToken, err := jwtManager.CreateToken(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value, ok := r.Context().Value(userIDKey).(uuid.UUID)
		if !ok {
			t.Fatalf("expected context value to be uuid.UUID")
		}

		if value != userID {
			t.Fatalf("expected userID %v, got %v", userID, value)
		}
	})

	handler := RequireAuth(jwtManager, next)

	request := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	request.Header.Set("Authorization", "Bearer "+testToken)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v", recorder.Code)
	}
}

func TestMiddlewareAuth_UserIDFromContext(t *testing.T) {
	userID, err := uuid.Parse("123e4567-e89b-12d3-a456-426614174000")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	ctx := context.WithValue(context.Background(), userIDKey, userID)

	result, ok := UserIDFromContext(ctx)

	if !ok {
		t.Fatalf("expected context value to uuid.UUID")
	}

	if result != userID {
		t.Fatalf("expected userID %v, got %v", userID, result)
	}
}

func TestMiddlewareAuth_UserIdFromContextMissing(t *testing.T) {
	ctx := context.Background()

	result, ok := UserIDFromContext(ctx)

	if ok {
		t.Fatalf("expected no context value")
	}

	if result != uuid.Nil() {
		t.Fatalf("expected no uuid value, got %v", result)
	}
}
