package main

import (
	"log"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/auth"
	"github.com/jwhite9387/taskflow-api/internal/config"
	"github.com/jwhite9387/taskflow-api/internal/handlers"
	"github.com/jwhite9387/taskflow-api/internal/repository"
	"github.com/jwhite9387/taskflow-api/internal/service"
)

func main() {
	repo := repository.NewMemoryUserRepository()
	secret, err := config.JWTSecret()
	if err != nil {
		log.Fatal(err)
	}
	tokenCreator := auth.NewJWTManager(secret)
	userService := service.NewUserService(repo, tokenCreator)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, userService)

	addr := ":" + config.Port()

	log.Fatal(http.ListenAndServe(addr, mux))
}
