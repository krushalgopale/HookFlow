package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func CreateEnvironment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter for tenantID
		tenantID := r.PathValue("tenant_id")

		// Get user id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check tenant ownership
		belongs, err := repository.TenantBelongsToUser(db, r.Context(), tenantID, userID)
		if err != nil {
			http.Error(w, "Failed to verify tenant ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var environment models.Environment

		// JSON Validation
		err = json.NewDecoder(r.Body).Decode(&environment)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if environment.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(models.EnvironmentErrRes{
				Error: "name is required",
			})
			return
		}

		// ID Generation
		environmentID := "env_" + uuid.New().String()

		// Database Operation
		err = repository.SaveEnvironment(db, r.Context(), environmentID, tenantID, environment.Name)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to create environment", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.EnvironmentResponse{
			Status: "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func GetEnvironment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		tenantID := r.PathValue("tenant_id")
		environmentID := r.PathValue("id")

		// User id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check Tenant Ownership
		belongs, err := repository.TenantBelongsToUser(db, r.Context(), tenantID, userID)
		if err != nil {
			http.Error(w, "Failed to verify tenant ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database operation
		environment, err := repository.GetEnvironmentByTenant(
			db,
			r.Context(),
			environmentID,
			tenantID,
		)
		// Error handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Environment not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch environment", http.StatusInternalServerError)
			return
		}

		// Server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(environment)
	}
}

func ListEnvronments(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		tenantID := r.PathValue("tenant_id")

		// Get user id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check Tenant Ownership
		belongs, err := repository.TenantBelongsToUser(db, r.Context(), tenantID, userID)
		if err != nil {
			http.Error(w, "Failed to verify tenant ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database Operation
		environments, err := repository.ListEnvironmentsByTenant(db, r.Context(), tenantID)
		// Error Handling
		if err != nil {
			http.Error(w, "Failed to fetch environments", http.StatusInternalServerError)
			return
		}

		response := models.EnvironmentListResponse{
			Environments: environments,
		}

		// Server Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func UpdateEnvironment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		tenantID := r.PathValue("tenant_id")
		environmentID := r.PathValue("id")

		// Get user Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check tenant ownership
		belongs, err := repository.TenantBelongsToUser(db, r.Context(), tenantID, userID)
		if err != nil {
			http.Error(w, "Failed to verify tenant ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var environment models.Environment

		// JSON validation
		err = json.NewDecoder(r.Body).Decode(&environment)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if environment.Name == "" {
			http.Error(w, "Environment name is required", http.StatusBadRequest)
			return
		}

		// Database operation
		err = repository.UpdateEnvironmentByTenant(
			db,
			r.Context(),
			environmentID,
			tenantID,
			environment.Name,
		)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to get updated environment", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.EnvironmentResponse{
			Status: "updated",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func DeleteEnvironment(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		tenantID := r.PathValue("tenant_id")
		environmentID := r.PathValue("id")

		// Get user Id from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check tenant ownership
		belongs, err := repository.TenantBelongsToUser(db, r.Context(), tenantID, userID)
		if err != nil {
			http.Error(w, "Failed to check tenant ownership", http.StatusInternalServerError)
			return
		}

		if !belongs {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Database operation
		err = repository.DeleteEnvironmentByTenant(db, r.Context(), environmentID, tenantID)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to delete environment", http.StatusInternalServerError)
			return
		}

		// Server response
		response := models.EnvironmentResponse{
			Status: "deleted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
