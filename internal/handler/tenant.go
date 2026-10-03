package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func CreateTenant(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tenant models.Tenant

		// JSON validation
		err := json.NewDecoder(r.Body).Decode(&tenant)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if tenant.Name == "" {
			w.Header().Set("Content Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(models.ErrResponse{
				Error: "organization name is required",
			})
			return
		}

		// Get User ID from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Tenant ID Generation
		tenant.ID = "org_" + uuid.New().String()

		// Database Operation
		err = repository.SaveTenant(db, r.Context(), tenant.ID, userID, tenant.Name)
		// Error Handling
		if err != nil {
			log.Println("Database Error", err)
			http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.TenantResponse{
			Status: "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func GetTenant(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		tenantID := r.PathValue("id")

		// Get User ID from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Database Operation
		tenant, err := repository.GetTenant(db, r.Context(), tenantID, userID)
		if err != nil {

			if err == pgx.ErrNoRows {
				http.Error(w, "Tenant not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to fetch Tenant", http.StatusInternalServerError)
			return
		}

		// Server Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(tenant)
	}
}

func ListTenants(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get User ID from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Database Operation
		tenants, err := repository.ListTenants(db, r.Context(), userID)
		// Error Handling
		if err != nil {
			http.Error(w, "Failed to fetch tenants", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.TenantListResponse{
			Tenants: tenants,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateTenant(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tenant models.Tenant

		// Get User ID from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Path parameter
		tenantID := r.PathValue("id")

		// JSON handling
		err := json.NewDecoder(r.Body).Decode(&tenant)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if tenant.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(models.ErrResponse{
				Error: "Tenant name is required",
			})
			return
		}

		// Database Operation
		err = repository.UpdateTenant(db, r.Context(), tenantID, userID, tenant.Name)
		// Error Handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Tenant not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to update tenant", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.TenantResponse{
			Status: "updated",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func DeleteTenant(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get User ID from request context
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Path parameter
		tenantID := r.PathValue("id")

		// Database Operation
		err := repository.DeleteTenant(db, r.Context(), tenantID, userID)
		// Error Handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "Tenant not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to delete tenant", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.TenantResponse{
			Status: "deleted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
