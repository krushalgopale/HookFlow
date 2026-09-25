package models

import "time"

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

type TenantListResponse struct {
	Tenants []Tenant `json:"tenants"`
} 
