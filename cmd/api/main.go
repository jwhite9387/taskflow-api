package main

import (
	"internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/health", handlers.Health)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
