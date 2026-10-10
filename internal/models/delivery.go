package models

import (
	"time"
)

type Delivery struct {
	ID              string     `json:"id"`
	EventID         string     `json:"event_id"`
	DestinationID   string     `json:"destination_id"`
	EventType       string     `json:"event_type"`
	DestinationName string     `json:"destination_name"`
	Status          string     `json:"status"`
	ResponseStatus  *int       `json:"response_status,omitempty"`
	ResponseBody    *string    `json:"response_body,omitempty"`
	Error           *string    `json:"error,omitempty"`
	AttemptCount    int        `json:"attempt_count"`
	MaxAttempts     int        `json:"max_attempts"`
	NextAttemptAt   *time.Time `json:"next_attempt_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type DeliveryAttempt struct {
	ID             string    `json:"id"`
	DeliveryID     string    `json:"delivery_id"`
	AttemptNumber  int       `json:"attempt_number"`
	Status         string    `json:"status"`
	ResponseStatus *int      `json:"response_status,omitempty"`
	ResponseBody   *string   `json:"response_body,omitempty"`
	Error          *string   `json:"error,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type DeliveryListItem struct {
	DestinationName string    `json:"destination_name"`
	EventType       string    `json:"event_type"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type DeliveryAttemptListItem struct {
	DestinationName string    `json:"destination_name"`
	AttemptNumber   int       `json:"attempt_number"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type DeliveriesResponse struct {
	Deliveries []DeliveryListItem `json:"deliveries"`
	Limit      int                `json:"limit,omitempty"`
	Offset     int                `json:"offset,omitempty"`
	Total      int                `json:"total,omitempty"`
}

type DeliveryAttemptsResponse struct {
	DeliveryAttempts []DeliveryAttemptListItem `json:"delivery_attempts"`
	Limit            int                       `json:"limit"`
	Offset           int                       `json:"offset"`
	Total            int                       `json:"total"`
}

type DeliveryResponse struct {
	Delivery []Delivery
}
