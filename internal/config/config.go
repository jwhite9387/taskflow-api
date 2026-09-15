package config

import (
	"errors"
	"os"
)

var ErrSecretNotSet = errors.New("secret not set")

func Port() string {
	val, ok := os.LookupEnv("PORT")
	if !ok {
		return "8080"
	}
	return val
}

func JWTSecret() (string, error) {
	val, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		return "", ErrSecretNotSet
	}
	return val, nil
}
