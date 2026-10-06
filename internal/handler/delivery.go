package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func GetEnvironmentDelivery(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		deliveryID := r.PathValue("id")
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

		// Database operation
		delivery, err := repository.GetDeliveryByEnvironment(db, r.Context(), deliveryID, environmentID)
		// Error handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "delivery not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to get delivery", http.StatusInternalServerError)
			return
		}

		// Server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(delivery)
	}
}

func ListEnvironmentDeliveries(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
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
			http.Error(w, "failed to check environment ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database operation
		deliveries, err := repository.ListDeliveriesByEnvironment(db, r.Context(), environmentID)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to fetch deliveries", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.DeliveriesResponse{
			Deliveries: deliveries,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}


func GetEventDelivery(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		deliveryID := r.PathValue("id")
		eventID := r.PathValue("evt_id")

		// Get User Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check environment ownership
		belongs, err := repository.EventBelongsToUser(db, r.Context(), eventID, userID)
		if err != nil {
			http.Error(w, "Failed to check event ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database operation
		delivery, err := repository.GetDeliveryByEnvironment(db, r.Context(), deliveryID, eventID)
		// Error handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "delivery not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to get delivery", http.StatusInternalServerError)
			return
		}

		// Server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(delivery)
	}
}


func ListEventDeliveries(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		eventID := r.PathValue("evt_id")

		// Get User Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check environment ownership
		belongs, err := repository.EventBelongsToUser(db, r.Context(), eventID, userID)
		if err != nil {
			http.Error(w, "failed to check event ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database operation
		deliveries, err := repository.ListDeliveriesByEvent(db, r.Context(), eventID)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to fetch deliveries", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.DeliveriesResponse{
			Deliveries: deliveries,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
