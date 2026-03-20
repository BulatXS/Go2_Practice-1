package http

import (
	"net/http"
	"os"

	"singularity.com/pr1/services/tasks/internal/client/authclient"
	"singularity.com/pr1/services/tasks/internal/service"
	"singularity.com/pr1/shared/middleware"
)

func NewRouter() http.Handler {
	svc := service.NewTaskService()

	authURL := os.Getenv("AUTH_BASE_URL")
	if authURL == "" {
		authURL = "http://localhost:8081"
	}

	client := authclient.New(authURL)
	handler := NewHandler(svc, client)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/tasks", handler.Tasks)
	mux.HandleFunc("/v1/tasks/", handler.TaskByID)

	return middleware.RequestID(middleware.Logging(mux))
}
