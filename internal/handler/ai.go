package handler

import (
	"encoding/json"
	"net/http"

	"github.com/krushalgopale/HookFlow/internal/ai"
	"github.com/krushalgopale/HookFlow/internal/models"
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
