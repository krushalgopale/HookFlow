package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
		
		// Tenant ID Generation
		tenant.ID = "ten" + uuid.New().String()
		
		// Database Operation
		err = repository.SaveTenant(db, r.Context(), tenant.ID, tenant.Name,)
		
		// Error Handling
		if err != nil {
			log.Println("Database Error", err)
			http.Error(w, "Failed to create tenant", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.TenantResponse{
			ID: tenant.ID,
			Status: "accepted",
		}

		w.Header().Set("Content Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		json.NewEncoder(w).Encode(response)

	}
}
