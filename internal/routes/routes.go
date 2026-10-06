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

	// Tenant
	mux.Handle("POST /tenant", middleware.Auth(handler.CreateTenant(db)))
	mux.Handle("GET /tenant/{id}", middleware.Auth(handler.GetTenant(db)))
	mux.Handle("GET /tenants", middleware.Auth(handler.ListTenants(db)))
	mux.Handle("PATCH /tenant/{id}", middleware.Auth(handler.UpdateTenant(db)))
	mux.Handle("DELETE /tenant/{id}", middleware.Auth(handler.DeleteTenant(db)))

	// Environment
	mux.Handle("POST /tenant/{tenant_id}/environment", middleware.Auth(handler.CreateEnvironment(db)))
	mux.Handle("GET /tenant/{tenant_id}/environment/{id}", middleware.Auth(handler.GetEnvironment(db)))
	mux.Handle("GET /tenant/{tenant_id}/environments", middleware.Auth(handler.ListEnvronments(db)))
	mux.Handle("PATCH /tenant/{tenant_id}/environment/{id}", middleware.Auth(handler.UpdateEnvironment(db)))
	mux.Handle("DELETE /tenant/{tenant_id}/environment/{id}", middleware.Auth(handler.DeleteEnvironment(db)))
	
	// API Key
	mux.Handle("POST /environment/{env_id}/api-key", middleware.Auth(handler.CreateAPIKey(db)))
	mux.Handle("GET /environment/{env_id}/api-key/{id}", middleware.Auth(handler.GetAPIKey(db)))
	mux.Handle("GET /environment/{env_id}/api-keys", middleware.Auth(handler.ListAPIKeys(db)))
	
	// Event
	mux.Handle("POST /event", middleware.ApiKeyAuth(db, handler.CreateEvent(db)))
	mux.Handle("GET /environment/{env_id}/event/{id}", middleware.Auth(handler.GetEvent(db)))
	mux.Handle("GET /environment/{env_id}/events", middleware.Auth(handler.ListEvents(db)))

	// Destination
	mux.Handle("POST /environment/{env_id}/destination", middleware.Auth(handler.CreateDestination(db)))
	mux.Handle("GET /environment/{env_id}/destination/{id}", middleware.Auth(handler.GetDestination(db)))
	mux.Handle("GET /environment/{env_id}/destinations", middleware.Auth(handler.ListDestinations(db)))
	mux.Handle("PATCH /environment/{env_id}/destination/{id}", middleware.Auth(handler.UpdateDestination(db)))
	mux.Handle("DELETE /environment/{env_id}/destination/{id}", middleware.Auth(handler.DeleteDestination(db)))
	
	// Delivery
	mux.Handle("GET /environment/{env_id}/deliveries", middleware.Auth(handler.ListEnvironmentDeliveries(db)))
	mux.Handle("GET /environment/{env_id}/delivery/{id}", middleware.Auth(handler.GetEnvironmentDelivery(db)))
	mux.Handle("GET /event/{evt_id}/deliveries", middleware.Auth(handler.ListEventDeliveries(db)))
	mux.Handle("GET /event/{evt_id}/delivery/{id}", middleware.Auth(handler.GetEventDelivery(db)))
	return mux
}
