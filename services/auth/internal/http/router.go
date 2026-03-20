package http

import (
	"net/http"

	"singularity.com/pr1/services/auth/internal/service"
	"singularity.com/pr1/shared/middleware"
)

func NewRouter() http.Handler {
	authService := service.NewAuthService()
	handler := NewHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/login", handler.Login)
	mux.HandleFunc("/v1/auth/verify", handler.Verify)

	return middleware.RequestID(middleware.Logging(mux))
}
