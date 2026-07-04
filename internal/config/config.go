package config

import (
	"os"
)

func Port() string {
	val, ok := os.LookupEnv("PORT")
	if !ok {
		return "8080"
	}
	return val
}
