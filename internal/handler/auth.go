package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/auth"
	"github.com/krushalgopale/HookFlow/internal/middleware"
	"github.com/krushalgopale/HookFlow/internal/models"
	"github.com/krushalgopale/HookFlow/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func Signup(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.SignUpRequest

		// JSON Handling
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.Email == "" && request.Password == "" {
			http.Error(w, "Email and Password are requried", http.StatusBadRequest)
			return
		}

		if request.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		if request.Password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		// Password hashing
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)
		// Error handling
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// ID generation
		userID := "usr_" + uuid.New().String()

		// Database Operation
		err = repository.SaveUser(db, r.Context(), userID, request.Email, string(passwordHash))
		// Error handling
		if err != nil {

			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				http.Error(w, "Email already exists", http.StatusConflict)
				return
			}
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		// Server Response
		response := models.AuthResponse{
			ID:     userID,
			Status: "created",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func Signin(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.SignInRequest

		// JSON Handling
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.Email == "" && request.Password == "" {
			http.Error(w, "Email and Password are required", http.StatusBadRequest)
			return
		}

		if request.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		if request.Password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		// Database Operation
		user, err := repository.GetUserByEmail(db, r.Context(), request.Email)
		// Error Handling
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Invalid email or password", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Failed to login", http.StatusInternalServerError)
			return
		}

		// Hashed Password Compare
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(request.Password),
		)
		// Error Handling
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Token Generation
		accessToken, err := auth.GenerateAccessToken(user.ID)
		if err != nil {
			http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
			return
		}

		// Set Authentication Cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    accessToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   15 * 60,
		})

		// Server Response
		response := map[string]string{
			"status": "signed in",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func Me(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIDKey).(string)

		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := repository.GetUserByID(db, r.Context(), userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}

		// Server Response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(user)
	}
}
