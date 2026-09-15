package models

import "time"

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
}

type EventResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type EventListResponse struct {
	Events []Event `json:"events"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
	Total  int     `json:"total"`
}
