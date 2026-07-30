package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, key model.APIKey) (model.APIKey, error) {
	const query = `
		INSERT INTO api_keys (id, name, key_hash, key_prefix, created_at, updated_at, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, key_hash, key_prefix, created_at, updated_at, revoked_at
	`

	return scanAPIKey(r.pool.QueryRow(ctx, query,
		key.ID,
		key.Name,
		key.KeyHash,
		key.KeyPrefix,
		key.CreatedAt,
		key.UpdatedAt,
		key.RevokedAt,
	))
}

func (r *PostgresRepository) FindByHash(ctx context.Context, hash string) (model.APIKey, error) {
	const query = `
		SELECT id, name, key_hash, key_prefix, created_at, updated_at, revoked_at
		FROM api_keys
		WHERE key_hash = $1
	`

	return scanAPIKey(r.pool.QueryRow(ctx, query, hash))
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (model.APIKey, error) {
	const query = `
		SELECT id, name, key_hash, key_prefix, created_at, updated_at, revoked_at
		FROM api_keys
		WHERE id = $1
	`

	return scanAPIKey(r.pool.QueryRow(ctx, query, id))
}

func (r *PostgresRepository) Revoke(ctx context.Context, id string, revokedAt time.Time) (model.APIKey, error) {
	const query = `
		UPDATE api_keys
		SET revoked_at = $2, updated_at = $2
		WHERE id = $1
		RETURNING id, name, key_hash, key_prefix, created_at, updated_at, revoked_at
	`

	return scanAPIKey(r.pool.QueryRow(ctx, query, id, revokedAt))
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanAPIKey(row rowScanner) (model.APIKey, error) {
	var key model.APIKey
	if err := row.Scan(
		&key.ID,
		&key.Name,
		&key.KeyHash,
		&key.KeyPrefix,
		&key.CreatedAt,
		&key.UpdatedAt,
		&key.RevokedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.APIKey{}, ErrNotFound
		}
		return model.APIKey{}, err
	}

	return key, nil
}
