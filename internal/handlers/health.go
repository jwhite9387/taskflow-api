package handlers

import (
	"io"
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "TaskFlow API is Running!\n")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
