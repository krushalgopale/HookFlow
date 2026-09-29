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
	mux.HandleFunc("POST /auth/signup", handler.Signup(db))
	mux.HandleFunc("POST /auth/signin", handler.Signin(db))
	mux.HandleFunc("POST /auth/signout", handler.Signout)
	mux.Handle("GET /auth/me", middleware.Auth(handler.Me(db)))
	mux.HandleFunc("POST /event", handler.CreateEvent(db))
	mux.HandleFunc("GET /event/{id}", handler.GetEvent(db))
	mux.HandleFunc("GET /events", handler.ListEvents(db))
	mux.HandleFunc("POST /tenant", handler.CreateTenant(db))
	mux.HandleFunc("GET /tenant/{id}", handler.GetTenant(db))
	mux.HandleFunc("GET /tenants", handler.ListTenants(db))
	mux.HandleFunc("PATCH /tenant/{id}", handler.UpdateTenant(db))
	mux.HandleFunc("DELETE /tenant/{id}", handler.DeleteTenant(db))
	return mux
}
