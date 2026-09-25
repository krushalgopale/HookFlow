package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/krushalgopale/HookFlow/internal/handler"
)

func Routes(db *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	// Registering Route using HandleFunc
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("POST /event", handler.CreateEvent(db))
	mux.HandleFunc("GET /event/{id}", handler.GetEvent(db))
	mux.HandleFunc("GET /events", handler.ListEvents(db))
	mux.HandleFunc("POST /tenant", handler.CreateTenant(db))
	mux.HandleFunc("GET /tenant/{id}", handler.GetTenant(db))
	mux.HandleFunc("GET /tenants", handler.ListTenants(db))
	mux.HandleFunc("PATCH /tenant/{id}", handler.UpdateTenant(db))
	return mux
}
