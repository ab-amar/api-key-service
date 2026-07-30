package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ab-amar/api-key-service/internal/model"
	"github.com/ab-amar/api-key-service/internal/repository"
)

type fakeRepository struct {
	byHash map[string]model.APIKey
	byID   map[string]model.APIKey
}

func newFakeRepository() *fakeRepository {
	return &fakeRepository{
		byHash: make(map[string]model.APIKey),
		byID:   make(map[string]model.APIKey),
	}
}

func (r *fakeRepository) Create(ctx context.Context, key model.APIKey) (model.APIKey, error) {
	r.byHash[key.KeyHash] = key
	r.byID[key.ID] = key
	return key, nil
}

func (r *fakeRepository) FindByHash(ctx context.Context, hash string) (model.APIKey, error) {
	key, ok := r.byHash[hash]
	if !ok {
		return model.APIKey{}, repository.ErrNotFound
	}
	return key, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, id string) (model.APIKey, error) {
	key, ok := r.byID[id]
	if !ok {
		return model.APIKey{}, repository.ErrNotFound
	}
	return key, nil
}

func (r *fakeRepository) Revoke(ctx context.Context, id string, revokedAt time.Time) (model.APIKey, error) {
	key, ok := r.byID[id]
	if !ok {
		return model.APIKey{}, repository.ErrNotFound
	}
	key.RevokedAt = &revokedAt
	key.UpdatedAt = revokedAt
	r.byID[id] = key
	r.byHash[key.KeyHash] = key
	return key, nil
}

func TestCreateReturnsPlainKeyAndStoresHash(t *testing.T) {
	repo := newFakeRepository()
	svc := NewAPIKeyService(repo)

	result, err := svc.Create(context.Background(), "billing")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if result.Plain == "" {
		t.Fatal("expected plain API key in create response")
	}
	if result.Key.KeyHash == result.Plain {
		t.Fatal("expected stored key hash to differ from plain key")
	}
	if result.Key.KeyPrefix == "" {
		t.Fatal("expected key prefix")
	}
}

func TestValidateRejectsUnknownKey(t *testing.T) {
	svc := NewAPIKeyService(newFakeRepository())

	_, err := svc.Validate(context.Background(), "missing")
	if !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("expected ErrInvalidKey, got %v", err)
	}
}

func TestValidateRejectsRevokedKey(t *testing.T) {
	repo := newFakeRepository()
	svc := NewAPIKeyService(repo)

	result, err := svc.Create(context.Background(), "billing")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if _, err := svc.Revoke(context.Background(), result.Key.ID); err != nil {
		t.Fatalf("Revoke returned error: %v", err)
	}

	_, err = svc.Validate(context.Background(), result.Plain)
	if !errors.Is(err, ErrKeyRevoked) {
		t.Fatalf("expected ErrKeyRevoked, got %v", err)
	}
}
