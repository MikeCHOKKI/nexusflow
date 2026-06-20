// ─── Router ────────────────────────────────────────────────────
package main

import (
	"net/http"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/handler"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

// newRouter assembles the full router with all routes and middleware.
func newRouter(clients *client.Clients, rdb *redis.Client) *mux.Router {
	r := mux.NewRouter()

	// ── Global middleware ──────────────────────────────────────
	r.Use(middleware.RateLimit(rdb))
	r.Use(corsMiddleware)
	r.Use(requestLogger)

	// ── Handlers ───────────────────────────────────────────────
	healthH := handler.NewHealthHandler(clients, clients.Conns)
	authH := handler.NewAuthHandler(clients)
	catalogH := handler.NewCatalogHandler(clients)
	orderH := handler.NewOrderHandler("service-order:50053")
	paymentH := handler.NewPaymentHandler(clients, rdb)
	dashboardH := handler.NewDashboardHandler(clients, clients.Conns, "service-order:50053")

	// ── Routes ─────────────────────────────────────────────────

	// Health — no auth required.
	r.HandleFunc("/health", healthH.Check).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/api/v1/health", healthH.Check).Methods(http.MethodGet, http.MethodOptions)

	// Auth — public.
	auth := r.PathPrefix("/api/v1/auth").Subrouter()
	auth.HandleFunc("/register", authH.Register).Methods(http.MethodPost)
	auth.HandleFunc("/login", authH.Login).Methods(http.MethodPost)

	// Auth-protected routes.
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.Auth(clients))

	// Auth me
	protected.HandleFunc("/auth/me", authH.Me).Methods(http.MethodGet)

	// Catalog
	protected.HandleFunc("/products", catalogH.ListProducts).Methods(http.MethodGet)
	protected.HandleFunc("/products/{id}", catalogH.GetProduct).Methods(http.MethodGet)
	protected.HandleFunc("/products", catalogH.CreateProduct).Methods(http.MethodPost)
	protected.HandleFunc("/products/{id}", catalogH.UpdateProduct).Methods(http.MethodPut)
	protected.HandleFunc("/products/{id}", catalogH.DeleteProduct).Methods(http.MethodDelete)

	// Orders
	protected.HandleFunc("/orders", orderH.CreateOrder).Methods(http.MethodPost)
	protected.HandleFunc("/orders", orderH.ListOrders).Methods(http.MethodGet)
	protected.HandleFunc("/orders/{id}", orderH.GetOrder).Methods(http.MethodGet)
	protected.HandleFunc("/orders/{id}/status", orderH.UpdateOrderStatus).Methods(http.MethodPut)

	// Payments
	protected.HandleFunc("/payments", paymentH.ProcessPayment).Methods(http.MethodPost)
	protected.HandleFunc("/payments/{id}", paymentH.GetPayment).Methods(http.MethodGet)
	protected.HandleFunc("/payments/{id}/refund", paymentH.RefundPayment).Methods(http.MethodPost)

	// Dashboard
	protected.HandleFunc("/dashboard/stats", dashboardH.Stats).Methods(http.MethodGet)

	return r
}

// corsMiddleware handles CORS headers.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requestLogger logs incoming HTTP requests.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simple passthrough — structured logging can be added via zerolog/prometheus.
		next.ServeHTTP(w, r)
	})
}
