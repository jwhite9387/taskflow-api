package main

import (
	"log"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/config"
	"github.com/jwhite9387/taskflow-api/internal/handlers"
	"github.com/jwhite9387/taskflow-api/internal/repository"
	"github.com/jwhite9387/taskflow-api/internal/service"
)

func main() {
	repo := repository.NewMemoryUserRepository()
	userService := service.NewUserService(repo)

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, userService)

	addr := ":" + config.Port()

	log.Fatal(http.ListenAndServe(addr, mux))
}
