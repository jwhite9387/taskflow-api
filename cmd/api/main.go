package main

import (
	"log"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/handlers"
)

func main() {

	mux := http.NewServeMux()

	handlers.RegisterRoutes(mux)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
