package models

import "time"

type Environment struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EnvironmentResponse struct {
	Status string `json:"status"`
}

type EnvironmentErrRes struct {
	Error string `json:"error"`
}

type EnvironmentListResponse struct {
	Environments []Environment `json:"tenants"`
}
