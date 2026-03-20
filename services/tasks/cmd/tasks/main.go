package main

import (
	"log"
	"net/http"
	"os"

	httpapi "singularity.com/pr1/services/tasks/internal/http"
)

func main() {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	addr := ":" + port

	log.Printf("tasks service started on %s", addr)
	if err := http.ListenAndServe(addr, httpapi.NewRouter()); err != nil {
		log.Fatalf("tasks service failed: %v", err)
	}
}
