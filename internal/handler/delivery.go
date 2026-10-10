package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

		// Database operation
		deliveries, total, err := repository.ListDeliveriesByEnvironment(db, r.Context(), environmentID, limit, offset)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to fetch deliveries", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.DeliveriesResponse{
			Deliveries: deliveries,
			Total: total,
			Limit: limit,
			Offset: offset,
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

		// Check event ownership
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
		delivery, err := repository.GetDeliveryByEvent(db, r.Context(), deliveryID, eventID)
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

		// Check event ownership
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

func ListDeliveryAttemptsByDelivery(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		deliveryID := r.PathValue("del_id")

		// Get User Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check delivery ownership
		belongs, err := repository.DeliveryBelongsToUser(db, r.Context(), deliveryID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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


		// Database operation
		deliveryAttempts, total, err := repository.ListDeliveryAttemptByDeliveryID(db, r.Context(), deliveryID, limit, offset)
		if err != nil {
			http.Error(w, "Failed to fetch delivery attempts", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.DeliveryAttemptsResponse{
			DeliveryAttempts: deliveryAttempts,
			Total: total,
			Limit: limit,
			Offset: offset,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
