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
	mux.HandleFunc("POST /", handler.CreateEvent(db))
	mux.HandleFunc("GET /event/{id}", handler.GetEvent(db))
	mux.HandleFunc("GET /events", handler.ListEvents(db))
	mux.HandleFunc("POST /tenants", handler.CreateTenant(db))
	return mux
}
