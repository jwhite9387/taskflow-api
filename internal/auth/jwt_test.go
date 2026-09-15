package auth

import (
	"testing"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

var testSecret = "testSecret"
var testUserID = "123e4567-e89b-12d3-a456-426614174000"

func TestJWT_CreateToken(t *testing.T) {
	manager := NewJWTManager(testSecret)
	parseId, err := uuid.Parse(testUserID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := manager.CreateToken(parseId)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token == "" {
		t.Fatalf("expected token, got empty string")
	}
}

func TestJWT_ValidateToken(t *testing.T) {
	manager := NewJWTManager(testSecret)
	parseID, err := uuid.Parse(testUserID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := manager.CreateToken(parseID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	validToken, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if validToken != parseID {
		t.Fatalf("expected user ID %v, got %v", parseID, validToken)
	}
}

func TestJWT_ValidateTokenWrongSecret(t *testing.T) {
	managerA := NewJWTManager(testSecret)
	parseID, err := uuid.Parse(testUserID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := managerA.CreateToken(parseID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	managerB := NewJWTManager("differentSecret")
	_, err = managerB.ValidateToken(token)
	if err == nil {
		t.Fatalf("expected validation to fail with wrong secret")
	}
}

func TestJWT_ValidateTokenExpired(t *testing.T) {
	manager := NewJWTManager(testSecret)

	claims := jwt.RegisteredClaims{
		Subject:   testUserID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = manager.ValidateToken(signedToken)
	if err == nil {
		t.Fatalf("expected validation to fail with expired time")
	}
}

func TestJWT_ValidateTokenWrongAlgorithm(t *testing.T) {
	manager := NewJWTManager(testSecret)

	claims := jwt.RegisteredClaims{
		Subject:   testUserID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS384,
		claims,
	)

	signedToken, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = manager.ValidateToken(signedToken)
	if err == nil {
		t.Fatalf("expected validation to fail with wrong algorithm")
	}
}
