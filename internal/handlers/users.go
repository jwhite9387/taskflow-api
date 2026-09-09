package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/service"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []string `json:"errors"`
}

type UserHandler struct {
	userService *service.UserService
}

func (r RegisterUserRequest) Validate() []string {
	var errs []string
	if r.Username == "" {
		errs = append(errs, "Username is required")
	}

	if r.Email == "" {
		errs = append(errs, "Email is required")
	}

	if r.Password == "" {
		errs = append(errs, "Password is required")
	} else if len(r.Password) < 8 {
		errs = append(errs, "Password must be at least 8 characters")
	}
	return errs
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userService.Register(user.Username, user.Email, user.Password); err != nil {
		var validationErr service.ValidationError

		if errors.As(err, &validationErr) {
			response := ValidationErrorResponse{
				Errors: validationErr.Errors,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				return
			}
			return
		}

		http.Error(w, "Unable to register user", http.StatusInternalServerError)
		return
	}

	response := RegisterUserResponse{
		Message: "User registered successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// TODO: Log this error once we add logging
		return
	}
}
