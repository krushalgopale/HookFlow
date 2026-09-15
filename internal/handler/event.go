package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func CreateEvent(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event models.Event

		// Json Validation
		err := json.NewDecoder(r.Body).Decode(&event)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if event.Data == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "data is required",
			})
			return
		}

		if event.Type == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(models.ErrorResponse{
				Error: "type is required",
			})
			return
		}

		// ID generation
		eventID := uuid.New().String()

		data, err := json.Marshal(event.Data)
		if err != nil {
			http.Error(w, "Failed to encode event data", http.StatusInternalServerError)
			return
		}

		// Database Operation
		err = repository.SaveEvent(db, r.Context(), eventID, event.Type, data)
		// Error Handling
		if err != nil {
			log.Println("Database error:", err)
			http.Error(w, "Failed to save event in database", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.EventResponse{
			ID:     eventID,
			Status: "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)

		// Logs
		fmt.Println(eventID)
		fmt.Println(event.Data)
		fmt.Println(event.Type)
	}
}

func GetEvent(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path parameter
		eventID := r.PathValue("id")

		// Database Operation
		event, err := repository.GetEvent(db, r.Context(), eventID)
		// Error Handling
		if err != nil {

			if err == pgx.ErrNoRows {
				http.Error(w, "Event not found", http.StatusNotFound)
				return
			}

			log.Println("Database error:", err)
			http.Error(w, "Failed to get event", http.StatusInternalServerError)
			return
		}

		// Server Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(event)
	}
}

func ListEvents(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query Parameter
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))

		// Defualt Limit Value
		if limit == 0 {
			limit = 10
		}

		// Validation of Query parameter
		if err != nil && r.URL.Query().Get("limit") != "" {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		if err != nil && r.URL.Query().Get("offset") != "" {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		if limit < 0 || offset < 0 {
			http.Error(w, "limit or offset cannot be negative", http.StatusBadRequest)
			return
		}

		if limit > 100 {
			http.Error(w, "limit cannot be greater than 100", http.StatusBadRequest)
			return
		}

		if offset > 1000000 {
			http.Error(w, "offset cannot be greater than 1000000", http.StatusBadRequest)
			return
		}

		// Database Operation
		events, total, err := repository.ListEvents(db, r.Context(), limit, offset)
		// Error Handling
		if err != nil {
			log.Println("Database error:", err)
			http.Error(w, "Failed to get events", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.EventListResponse{
			Events: events,
			Limit:  limit,
			Offset: offset,
			Total:  total,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
