package models

import "time"

type Destination struct {
	ID            string    `json:"id"`
	EnvironmentID string    `json:"environment_id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DestinationResponse struct {
	Status string `json:"status"`
}

type ListDestinationResponse struct {
	Destinations []Destination `json:"destinations"`
}

type CreateDestinationRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type UpdateDestinationRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
