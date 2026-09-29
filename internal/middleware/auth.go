package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication Cookie
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Verify JWT
		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (any, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("unexpected signing method")
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract and check JWT claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token clams", http.StatusUnauthorized)
			return
		}

		// Extract and check user ID
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Invalid user ID", http.StatusUnauthorized)
			return
		}

		// Store User ID in Request Context
		ctx := context.WithValue(
			r.Context(),
			UserIDKey,
			userID,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
