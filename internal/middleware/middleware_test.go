package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ab-amar/api-key-service/internal/metrics"
	"github.com/ab-amar/api-key-service/internal/model"
)

type fakeValidator struct {
	key model.APIKey
	err error
}

func (v fakeValidator) Validate(ctx context.Context, plain string) (model.APIKey, error) {
	return v.key, v.err
}

func TestAuthenticateStoresAPIKeyInContext(t *testing.T) {
	key := model.APIKey{ID: "key_1"}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := APIKeyFromContext(r.Context())
		if !ok {
			t.Fatal("expected api key in context")
		}
		if got.ID != key.ID {
			t.Fatalf("expected key id %q, got %q", key.ID, got.ID)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	handler := Authenticate(fakeValidator{key: key})(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer ak_test")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

func TestTokenBucketReturns429WhenExhausted(t *testing.T) {
	limiter := NewTokenBucketLimiter(1, 0, metrics.New())
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	handler := limiter.Middleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), apiKeyKey, model.APIKey{ID: "key_1"})
	req = req.WithContext(ctx)

	first := httptest.NewRecorder()
	handler.ServeHTTP(first, req)
	if first.Code != http.StatusNoContent {
		t.Fatalf("expected first request to pass, got %d", first.Code)
	}

	second := httptest.NewRecorder()
	handler.ServeHTTP(second, req)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request to be rate limited, got %d", second.Code)
	}
}
