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

func CreateAPIKey(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// ID and Key Generation
		apiKeyID := "key_" + uuid.New().String()
		apiKey := "hf_" + uuid.New().String()

		// Database operation
		err = repository.SaveAPIKey(db, r.Context(), apiKeyID, environmentID, apiKey)
		if err != nil {
			http.Error(w, "Failed to create API key", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.APIKeyResponse{
			APIKey: apiKey,
			Status: "accepted",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func GetAPIKey(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Path parameter
		apiKeyID := r.PathValue("id")
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
		apiKey, err := repository.GetAPIKeyByID(db, r.Context(), apiKeyID, environmentID)
		// Error handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "API Key not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to get API key", http.StatusInternalServerError)
			return
		}

		// Server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(apiKey)
	}
}

func ListAPIKeys(db *pgxpool.Pool) http.HandlerFunc {
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

		ApiKeys, err := repository.ListAPIKeysByEnvironment(db, r.Context(), environmentID)
		if err != nil {
			http.Error(w, "Failed to fetch api keys", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.APIKeyListResponse{
			APIKeys: ApiKeys,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}
