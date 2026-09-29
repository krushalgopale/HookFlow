package models

import "time"

type Tenant struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantResponse struct {
	Status string `json:"status"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type TenantListResponse struct {
	Tenants []Tenant `json:"tenants"`
}
