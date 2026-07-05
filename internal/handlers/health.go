package handlers

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := HealthResponse{
		Status:  "ok",
		Service: "TaskFlow API",
		Version: "0.1.0",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
}
