package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/ab-amar/api-key-service/internal/repository"
)

var (
	ErrInvalidName = errors.New("key name is required")
	ErrInvalidKey  = errors.New("api key is invalid")
	ErrKeyRevoked  = errors.New("api key is revoked")
)

type APIKeyService struct {
	repo repository.APIKeyRepository
}

type CreateAPIKeyResult struct {
	Key   model.APIKey
	Plain string
}

func NewAPIKeyService(repo repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) Create(ctx context.Context, name string) (CreateAPIKeyResult, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return CreateAPIKeyResult{}, ErrInvalidName
	}

	plain, err := generateKey()
	if err != nil {
		return CreateAPIKeyResult{}, err
	}

	now := time.Now().UTC()
	key := model.APIKey{
		ID:        newID(now),
		Name:      name,
		KeyHash:   HashKey(plain),
		KeyPrefix: plain[:12],
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repo.Create(ctx, key)
	if err != nil {
		return CreateAPIKeyResult{}, err
	}

	return CreateAPIKeyResult{Key: created, Plain: plain}, nil
}

func (s *APIKeyService) Validate(ctx context.Context, plain string) (model.APIKey, error) {
	plain = strings.TrimSpace(plain)
	if plain == "" {
		return model.APIKey{}, ErrInvalidKey
	}

	key, err := s.repo.FindByHash(ctx, HashKey(plain))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.APIKey{}, ErrInvalidKey
		}
		return model.APIKey{}, err
	}

	if key.IsRevoked() {
		return model.APIKey{}, ErrKeyRevoked
	}

	return key, nil
}

func (s *APIKeyService) Revoke(ctx context.Context, id string) (model.APIKey, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return model.APIKey{}, repository.ErrNotFound
	}

	return s.repo.Revoke(ctx, id, time.Now().UTC())
}

func HashKey(plain string) string {
	sum := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(sum[:])
}

func generateKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate api key: %w", err)
	}

	return "ak_" + hex.EncodeToString(bytes), nil
}

func newID(now time.Time) string {
	random := make([]byte, 8)
	if _, err := rand.Read(random); err != nil {
		return fmt.Sprintf("key_%d", now.UnixNano())
	}

	return fmt.Sprintf("key_%d_%s", now.UnixNano(), hex.EncodeToString(random))
}
