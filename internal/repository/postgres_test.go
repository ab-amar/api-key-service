package repository

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRepositoryCreateFindAndRevoke(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL is not set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS api_keys (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			key_hash TEXT NOT NULL UNIQUE,
			key_prefix TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			revoked_at TIMESTAMPTZ
		)
	`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, "DELETE FROM api_keys WHERE id = $1", "test_key_1")
	})

	repo := NewPostgresRepository(pool)
	now := time.Now().UTC()
	key := model.APIKey{
		ID:        "test_key_1",
		Name:      "integration",
		KeyHash:   "test_hash_1",
		KeyPrefix: "ak_test",
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := repo.Create(ctx, key)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != key.ID {
		t.Fatalf("expected id %q, got %q", key.ID, created.ID)
	}

	found, err := repo.FindByHash(ctx, key.KeyHash)
	if err != nil {
		t.Fatalf("FindByHash returned error: %v", err)
	}
	if found.Name != key.Name {
		t.Fatalf("expected name %q, got %q", key.Name, found.Name)
	}

	revoked, err := repo.Revoke(ctx, key.ID, now.Add(time.Second))
	if err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}
	if revoked.RevokedAt == nil {
		t.Fatal("expected revoked_at to be set")
	}
}
