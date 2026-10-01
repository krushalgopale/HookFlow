package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/handler"
	"github.com/krushalgopale/HookFlow/internal/middleware"
)

func Routes(db *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Registering Route using HandleFunc/Handle
	mux.HandleFunc("GET /health", handler.Health)

	// Auth
	mux.HandleFunc("POST /auth/signup", handler.Signup(db))
	mux.HandleFunc("POST /auth/signin", handler.Signin(db))
	mux.HandleFunc("POST /auth/signout", handler.Signout)
	mux.Handle("GET /auth/me", middleware.Auth(handler.Me(db)))

	// Event
	mux.HandleFunc("POST /event", handler.CreateEvent(db))
	mux.HandleFunc("GET /event/{id}", handler.GetEvent(db))
	mux.HandleFunc("GET /events", handler.ListEvents(db))

	// Tenant
	mux.Handle("POST /tenant", middleware.Auth(handler.CreateTenant(db)))
	mux.Handle("GET /tenant/{id}", middleware.Auth(handler.GetTenant(db)))
	mux.Handle("GET /tenants", middleware.Auth(handler.ListTenants(db)))
	mux.Handle("PATCH /tenant/{id}", middleware.Auth(handler.UpdateTenant(db)))
	mux.Handle("DELETE /tenant/{id}", middleware.Auth(handler.DeleteTenant(db)))
	
	// Environment
	mux.Handle("POST /tenant/{id}/environment", middleware.Auth(handler.CreateEnvironment(db)))

	return mux
}
