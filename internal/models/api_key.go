package models

import (
	"time"
)

type APIKey struct {
	ID            string    `json:"id"`
	EnvironmentID string    `json:"environment_id"`
	Key           string    `json:"key"`
	CreatedAt     time.Time `json:"created_at"`
}

type APIKeyResponse struct {
	APIKey string `json:"api-key"`
	Status string `json:"status"`
}

type ErrRes struct {
	Error string `json:"error"`
}
