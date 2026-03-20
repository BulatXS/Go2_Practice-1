package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"singularity.com/pr1/services/tasks/internal/client/authclient"
	"singularity.com/pr1/services/tasks/internal/service"
	"singularity.com/pr1/shared/middleware"
)

type Handler struct {
	svc        *service.TaskService
	authClient *authclient.Client
}

func NewHandler(svc *service.TaskService, authClient *authclient.Client) *Handler {
	return &Handler{
		svc:        svc,
		authClient: authClient,
	}
}

func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	reqID := middleware.GetRequestID(r.Context())

	ok, status, err := h.authClient.Verify(r.Context(), token, reqID)
	if err != nil {
		http.Error(w, "auth service unavailable", http.StatusBadGateway)
		return false
	}

	if !ok {
		if status == http.StatusUnauthorized || status == http.StatusForbidden {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return false
		}

		http.Error(w, "auth error", http.StatusBadGateway)
		return false
	}

	return true
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.CreateTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskByID(w http.ResponseWriter, r *http.Request) {
	if !h.authorize(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPatch:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	task := h.svc.Create(req.Title, req.Description, req.DueDate)
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.svc.GetAll())
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")

	task, ok := h.svc.GetByID(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		DueDate     *string `json:"due_date"`
		Done        *bool   `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	task, ok := h.svc.Update(id, req.Title, req.Description, req.DueDate, req.Done)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")

	if ok := h.svc.Delete(id); !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
