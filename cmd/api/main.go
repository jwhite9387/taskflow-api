package main

import (
	"log"
	"net/http"

	"github.com/jwhite9387/taskflow-api/internal/config"
	"github.com/jwhite9387/taskflow-api/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	addr := ":" + config.Port()

	log.Fatal(http.ListenAndServe(addr, mux))
}
