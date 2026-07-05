package handlers

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/users/register", RegisterUser)
}
