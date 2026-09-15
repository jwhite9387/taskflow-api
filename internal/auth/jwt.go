package auth

import (
	"errors"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secret: secret,
	}
}

func (m *JWTManager) CreateToken(userID uuid.UUID) (string, error) {
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString([]byte(m.secret))

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("token is not using HS256")
			}
			return []byte(m.secret), nil
		},
	)

	if err != nil {
		return uuid.Nil(), err
	}

	if !token.Valid {
		return uuid.Nil(), errors.New("invalid token")
	}

	sub, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil(), err
	}
	return sub, nil
}
