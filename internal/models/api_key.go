package models

import (
	"time"
)

type APIKey struct {
	ID            string     `json:"id"`
	EnvironmentID string     `json:"environment_id"`
	Name          string     `json:"name"`
	Key           string     `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
}

type APIKeyResponse struct {
	APIKey string `json:"api-key"`
	Status string `json:"status"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type ErrRes struct {
	Error string `json:"error"`
}

type CreateAPIKeyRequest struct {
	Name       string     `json:"name"`
	Expiration string     `json:"expiration"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

type APIKeyListResponse struct {
	APIKeys []APIKey `json:"api_keys"`
}
