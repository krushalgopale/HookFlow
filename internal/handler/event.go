package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
	"github.com/krushalgopale/HookFlow/internal/worker"
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

		// Get environment Id added by api key middleware
		environmentID, ok := r.Context().Value(middleware.EnvironmentIDKey).(string)
		if !ok {
			http.Error(w, "environment not found", http.StatusUnauthorized)
			return
		}

		// ID generation
		eventID := "evt_" + uuid.New().String()

		data, err := json.Marshal(event.Data)
		if err != nil {
			http.Error(w, "Failed to encode event data", http.StatusInternalServerError)
			return
		}

		// Database Operation
		err = repository.SaveEvent(db, r.Context(), eventID, environmentID, event.Type, data)
		// Error Handling
		if err != nil {
			http.Error(w, "Failed to save event in database", http.StatusInternalServerError)
			return
		}

		// Get destinations for the environment
		destinations, err := repository.ListDestinations(db, r.Context(), environmentID)
		if err != nil {
			http.Error(w, "Failed to fetch destinations", http.StatusInternalServerError)
			return
		}

		// Create a delivery for each destination
		for _, destination := range destinations {
			deliveryID := "del_" + uuid.New().String()

			err := repository.SaveDelivery(db, r.Context(), deliveryID, eventID, destination.ID)
			if err != nil {
				http.Error(w, "Failed to create delivery", http.StatusInternalServerError)
				return
			}

			go worker.ExecuteDelivery(
				db,
				context.Background(),
				deliveryID,
				1,
				)
		}

		// Server Response
		response := models.EventResponse{
			ID:     eventID,
			Status: "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func GetEvent(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// path parameter
		eventID := r.PathValue("id")
		environmentID := r.PathValue("env_id")

		// Get User Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check environment ownership
		belongs, err := repository.EnvironmentBelongsToUser(db, r.Context(), environmentID, userID)
		if err != nil {
			http.Error(w, "Failed to check environment ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database Operation
		event, err := repository.GetEvent(db, r.Context(), eventID, environmentID)
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
		// Path parameter
		environmentID := r.PathValue("env_id")

		// Get user Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check environment ownership
		belongs, err := repository.EnvironmentBelongsToUser(db, r.Context(), environmentID, userID)
		if err != nil {
			http.Error(w, "Failed to check environment ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Query Parameter
		limit, limitErr := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, offsetErr := strconv.Atoi(r.URL.Query().Get("offset"))

		// Defualt Limit Value
		if limit == 0 {
			limit = 10
		}

		// Validation of Query parameter
		if limitErr != nil && r.URL.Query().Get("limit") != "" {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		if offsetErr != nil && r.URL.Query().Get("offset") != "" {
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
		events, total, err := repository.ListEvents(db, r.Context(), environmentID, limit, offset)
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
