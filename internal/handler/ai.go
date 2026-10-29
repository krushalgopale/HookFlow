package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/ai"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
)

func AnalyzeFailure(client *ai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.FailureRequest

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if input.StatusCode < 100 || input.StatusCode > 599 ||
			input.Attempt < 1 ||
			len(input.ResponseBody) > 5000 {
			http.Error(w, "Invalid failure details", http.StatusBadRequest)
			return
		}

		result, err := client.Analyze(r.Context(), input)
		if err != nil {
			http.Error(w, "AI analysis service unavailable", http.StatusBadGateway)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}


func InvestigateDelivery(db *pgxpool.Pool, client *ai.Client) http.HandlerFunc {
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

		// Database operation
		attempts, err := repository.GetDeliveryAttemptsForInvestigation(db, r.Context(), deliveryID)
		// Error handling
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "delivery attempts not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to get delivery attempts", http.StatusInternalServerError)
			return
		}

		// AI operation
		input := models.InvestigationRequest{
			Delivery: delivery,
			Attempts: attempts,
		}

		result, err := client.Investigate(r.Context(), input)
		// Error handling
		if err != nil {
			http.Error(w, "AI investigation service unavailable", http.StatusBadGateway)
			return
		}

		// Server response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(result)
	}
}
