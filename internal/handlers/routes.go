package handlers

import (
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, userService *service.UserService) {
	userHandler := &UserHandler{
		userService: userService,
	}

	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/users/register", userHandler.RegisterUser)
	mux.HandleFunc("/users/login", userHandler.LoginUser)
}
