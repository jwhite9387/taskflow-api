package handlers

import (
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/auth"
	"github.com/jwhite9387/taskflow-api/internal/middleware"
	"github.com/jwhite9387/taskflow-api/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, userService *service.UserService, jwtManager *auth.JWTManager) {
	userHandler := &UserHandler{
		userService: userService,
	}

	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/users/register", userHandler.RegisterUser)
	mux.HandleFunc("/users/login", userHandler.LoginUser)
	mux.Handle("/users/me", middleware.RequireAuth(jwtManager, http.HandlerFunc(userHandler.GetCurrentUser)))
}
