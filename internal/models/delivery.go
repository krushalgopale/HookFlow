package models

import "time"

type Delivery struct {
	ID             string    `json:"id"`
	EventID        string    `json:"event_id"`
	DestinationID  string    `json:"destination_id"`
	EventType      string    `json:"event_type"`
	Status         string    `json:"status"`
	ResponseStatus *int      `json:"response_status,omitempty"`
	ResponseBody   *string   `json:"response_body,omitempty"`
	Error          *string   `json:"error,omitempty"`
	AttemptCount   int       `json:"attempt_count"`
	MaxAttempts    int       `json:"max_attempts"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DeliveryListItem struct {
	EventType string    `json:"event_type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DeliveriesResponse struct {
	Deliveries []DeliveryListItem
}

type DeliveryResponse struct {
	Delivery []Delivery
}
