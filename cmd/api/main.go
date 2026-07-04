package main

import (
	"log"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/handlers"
)

func main() {
	http.HandleFunc("/health", handlers.Health)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
