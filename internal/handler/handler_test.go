package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ab-amar/api-key-service/internal/metrics"
	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/ab-amar/api-key-service/internal/repository"
	"github.com/ab-amar/api-key-service/internal/service"
	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	createResult service.CreateAPIKeyResult
	createErr    error
	validateKey  model.APIKey
	validateErr  error
	revokeKey    model.APIKey
	revokeErr    error
}

func (s *fakeService) Create(ctx context.Context, name string) (service.CreateAPIKeyResult, error) {
	return s.createResult, s.createErr
}

func (s *fakeService) Validate(ctx context.Context, plain string) (model.APIKey, error) {
	return s.validateKey, s.validateErr
}

func (s *fakeService) Revoke(ctx context.Context, id string) (model.APIKey, error) {
	return s.revokeKey, s.revokeErr
}

type fakeReadiness struct {
	err error
}

func (r fakeReadiness) Ping(ctx context.Context) error {
	return r.err
}

func newTestRouter(service APIKeyService, readiness ReadinessChecker) http.Handler {
	router := chi.NewRouter()
	h := New(service, metrics.New(), readiness)
	h.RegisterRoutes(router)
	return router
}

func TestHealthReturnsOK(t *testing.T) {
	router := newTestRouter(&fakeService{}, fakeReadiness{})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("expected application/json, got %q", rec.Header().Get("Content-Type"))
	}
}

func TestReadyReturnsUnavailableWhenDependencyFails(t *testing.T) {
	router := newTestRouter(&fakeService{}, fakeReadiness{err: errors.New("db down")})
	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}
}

func TestCreateKeyReturnsPlainKeyOnce(t *testing.T) {
	now := time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	router := newTestRouter(&fakeService{
		createResult: service.CreateAPIKeyResult{
			Key: model.APIKey{
				ID:        "key_1",
				Name:      "billing",
				KeyPrefix: "ak_123456789",
				CreatedAt: now,
			},
			Plain: "ak_123456789abcdef",
		},
	}, fakeReadiness{})

	req := httptest.NewRequest(http.MethodPost, "/keys", bytes.NewBufferString(`{"name":"billing"}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["api_key"] != "ak_123456789abcdef" {
		t.Fatalf("expected plain api key in response, got %q", body["api_key"])
	}
}

func TestRevokeNotFoundReturns404(t *testing.T) {
	router := newTestRouter(&fakeService{revokeErr: repository.ErrNotFound}, fakeReadiness{})
	req := httptest.NewRequest(http.MethodPost, "/keys/missing/revoke", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
