package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ab-amar/api-key-service/internal/metrics"
	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/ab-amar/api-key-service/internal/service"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	apiKeyKey    contextKey = "api_key"
)

type APIKeyValidator interface {
	Validate(ctx context.Context, plain string) (model.APIKey, error)
}

func Chain(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	return handler
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}

		w.Header().Set("X-Request-ID", requestID)
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestLogger(logger *slog.Logger, m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			m.IncRequests()

			next.ServeHTTP(w, r)

			logger.Info("request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", RequestIDFromContext(r.Context()),
			)
		})
	}
}

func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered", "error", recovered, "request_id", RequestIDFromContext(r.Context()))
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Authenticate(validator APIKeyValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeText(w, http.StatusUnauthorized, "missing api key")
				return
			}

			plain := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if plain == "" {
				writeText(w, http.StatusUnauthorized, "missing api key")
				return
			}

			key, err := validator.Validate(r.Context(), plain)
			if err != nil {
				status := http.StatusUnauthorized
				if errors.Is(err, service.ErrKeyRevoked) {
					status = http.StatusForbidden
				}
				writeText(w, status, "invalid api key")
				return
			}

			ctx := context.WithValue(r.Context(), apiKeyKey, key)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type TokenBucketLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	capacity float64
	refill   float64
	metrics  *metrics.Metrics
}

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

func NewTokenBucketLimiter(capacity int, refillPerSecond int, m *metrics.Metrics) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:  make(map[string]*bucket),
		capacity: float64(capacity),
		refill:   float64(refillPerSecond),
		metrics:  m,
	}
}

func (l *TokenBucketLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, ok := APIKeyFromContext(r.Context())
		if !ok {
			writeText(w, http.StatusUnauthorized, "missing api key")
			return
		}

		if !l.allow(key.ID) {
			l.metrics.IncRateLimited()
			w.Header().Set("Retry-After", "1")
			writeText(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (l *TokenBucketLimiter) allow(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[id]
	if !ok {
		l.buckets[id] = &bucket{tokens: l.capacity - 1, lastRefill: now}
		return true
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens = min(l.capacity, b.tokens+elapsed*l.refill)
	b.lastRefill = now

	if b.tokens < 1 {
		return false
	}

	b.tokens--
	return true
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey).(string)
	return requestID
}

func APIKeyFromContext(ctx context.Context) (model.APIKey, bool) {
	key, ok := ctx.Value(apiKeyKey).(model.APIKey)
	return key, ok
}

func newRequestID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(bytes)
}

func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body + "\n"))
}
