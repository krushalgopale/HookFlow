package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/auth"
)

type ctxKey string

const EnvironmentIDKey ctxKey = "environmentID"

func ApiKeyAuth(db *pgxpool.Pool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Read the api key from the request header
		apiKey := r.Header.Get("X-API-Key")

		if apiKey == "" {
			http.Error(w, "API key is required", http.StatusUnauthorized)
			return
		}

		// Verify the API key and check its expiration
		result, err := auth.VerifyAPIKey(db, r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "Request canceled", http.StatusRequestTimeout)
				return
			}

			if err.Error() == "api key has expired" {
				http.Error(w, "api key has expired", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Invaild API key", http.StatusUnauthorized)
			return
		}

		// Store the autheticated environmentId in the request context
		ctx := context.WithValue(
			r.Context(),
			EnvironmentIDKey,
			result.EnvironmentID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
