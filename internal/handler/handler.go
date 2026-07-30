package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ab-amar/api-key-service/internal/metrics"
	"github.com/ab-amar/api-key-service/internal/middleware"
	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/ab-amar/api-key-service/internal/repository"
	"github.com/ab-amar/api-key-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type APIKeyService interface {
	Create(ctx context.Context, name string) (service.CreateAPIKeyResult, error)
	Validate(ctx context.Context, plain string) (model.APIKey, error)
	Revoke(ctx context.Context, id string) (model.APIKey, error)
}

type ReadinessChecker interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	service   APIKeyService
	metrics   *metrics.Metrics
	readiness ReadinessChecker
}

func New(service APIKeyService, metrics *metrics.Metrics, readiness ReadinessChecker) *Handler {
	return &Handler{
		service:   service,
		metrics:   metrics,
		readiness: readiness,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.root)
	router.Get("/health", h.health)
	router.Get("/ready", h.ready)
	router.Get("/metrics", h.metricsSnapshot)
	router.Post("/keys", h.createKey)
	router.Post("/auth/validate", h.validateKey)
	router.Post("/keys/{id}/revoke", h.revokeKey)
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"service": "api-key-service",
		"health":  "/health",
		"ready":   "/ready",
	})
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if h.readiness != nil {
		if err := h.readiness.Ping(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "service is not ready")
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

type createKeyRequest struct {
	Name string `json:"name"`
}

type createKeyResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	KeyPrefix string `json:"key_prefix"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) createKey(w http.ResponseWriter, r *http.Request) {
	var req createKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	result, err := h.service.Create(r.Context(), req.Name)
	if err != nil {
		if errors.Is(err, service.ErrInvalidName) {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create api key")
		return
	}

	h.metrics.IncAPIKeysCreated()
	writeJSON(w, http.StatusCreated, createKeyResponse{
		ID:        result.Key.ID,
		Name:      result.Key.Name,
		APIKey:    result.Plain,
		KeyPrefix: result.Key.KeyPrefix,
		CreatedAt: result.Key.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type validateKeyRequest struct {
	APIKey string `json:"api_key"`
}

type validateKeyResponse struct {
	Valid     bool   `json:"valid"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	KeyPrefix string `json:"key_prefix,omitempty"`
}

func (h *Handler) validateKey(w http.ResponseWriter, r *http.Request) {
	var req validateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	key, err := h.service.Validate(r.Context(), req.APIKey)
	if err != nil {
		if errors.Is(err, service.ErrInvalidKey) || errors.Is(err, service.ErrKeyRevoked) {
			writeJSON(w, http.StatusUnauthorized, validateKeyResponse{Valid: false})
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to validate api key")
		return
	}

	h.metrics.IncValidations()
	writeJSON(w, http.StatusOK, validateKeyResponse{
		Valid:     true,
		ID:        key.ID,
		Name:      key.Name,
		KeyPrefix: key.KeyPrefix,
	})
}

func (h *Handler) revokeKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	key, err := h.service.Revoke(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "api key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to revoke api key")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id":     key.ID,
		"status": "revoked",
	})
}

func (h *Handler) metricsSnapshot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.metrics.Snapshot())
}

func ProtectedResource(w http.ResponseWriter, r *http.Request) {
	key, _ := middleware.APIKeyFromContext(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "authenticated request accepted",
		"key_id":  key.ID,
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
