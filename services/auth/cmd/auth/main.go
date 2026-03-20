package main

import (
	"log"
	"net/http"
	"os"

	httpapi "singularity.com/pr1/services/auth/internal/http"
)

func main() {
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	addr := ":" + port
	handler := httpapi.NewRouter()

	log.Printf("auth service started on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("auth service failed: %v", err)
	}
}
