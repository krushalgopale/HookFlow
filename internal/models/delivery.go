package models

import (
	"time"
)

type Delivery struct {
	ID             string     `json:"id"`
	EventID        string     `json:"event_id"`
	DestinationID  string     `json:"destination_id"`
	EventType      string     `json:"event_type"`
	Status         string     `json:"status"`
	ResponseStatus *int       `json:"response_status,omitempty"`
	ResponseBody   *string    `json:"response_body,omitempty"`
	Error          *string    `json:"error,omitempty"`
	AttemptCount   int        `json:"attempt_count"`
	MaxAttempts    int        `json:"max_attempts"`
	NextAttemptAt  *time.Time `json:"next_attempt_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type DeliveryAttempt struct {
	ID             string    `json:"id"`
	DeliveryID     string    `json:"delivery_id"`
	AttemptNmuber  int       `json:"attempt_number"`
	Status         string    `json:"status"`
	ResponseStatus *int      `json:"response_status,omitempty"`
	ResponseBody   *string   `json:"response_body,omitempty"`
	Error          *string   `json:"error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
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
