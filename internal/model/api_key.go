package model

import "time"

type APIKey struct {
	ID        string
	Name      string
	KeyHash   string
	KeyPrefix string
	CreatedAt time.Time
	UpdatedAt time.Time
	RevokedAt *time.Time
}

func (k APIKey) IsRevoked() bool {
	return k.RevokedAt != nil
}
