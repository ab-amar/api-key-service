package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ab-amar/api-key-service/internal/model"
)

var ErrNotFound = errors.New("api key not found")

type APIKeyRepository interface {
	Create(ctx context.Context, key model.APIKey) (model.APIKey, error)
	FindByHash(ctx context.Context, hash string) (model.APIKey, error)
	FindByID(ctx context.Context, id string) (model.APIKey, error)
	Revoke(ctx context.Context, id string, revokedAt time.Time) (model.APIKey, error)
}
