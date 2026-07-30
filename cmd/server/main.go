package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ab-amar/api-key-service/internal/config"
	"github.com/ab-amar/api-key-service/internal/handler"
	"github.com/ab-amar/api-key-service/internal/metrics"
	"github.com/ab-amar/api-key-service/internal/middleware"
	"github.com/ab-amar/api-key-service/internal/repository"
	"github.com/ab-amar/api-key-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping postgres: %v", err)
	}

	repo := repository.NewPostgresRepository(pool)
	apiKeyService := service.NewAPIKeyService(repo)
	appMetrics := metrics.New()
	router := setupRouter(apiKeyService, appMetrics, pool, logger)
	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("server listening", "addr", cfg.Addr())
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	logger.Info("server stopped")
}

func setupRouter(apiKeyService *service.APIKeyService, appMetrics *metrics.Metrics, readiness handler.ReadinessChecker, logger *slog.Logger) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RequestLogger(logger, appMetrics))
	router.Use(middleware.Recover(logger))

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	apiHandler := handler.New(apiKeyService, appMetrics, readiness)
	limiter := middleware.NewTokenBucketLimiter(5, 1, appMetrics)
	apiHandler.RegisterRoutes(router)
	router.With(middleware.Authenticate(apiKeyService), limiter.Middleware).Get("/v1/protected", handler.ProtectedResource)

	return router
}
