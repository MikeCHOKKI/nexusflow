// ─── API Gateway — Entry Point ─────────────────────────────────
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/config"
	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/middleware"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[gateway] starting nexusflow api gateway...")

	// ── Configuration ──────────────────────────────────────────
	cfg := config.Load()

	// ── Redis ──────────────────────────────────────────────────
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rdb, err := middleware.NewRedisClient(ctx, cfg.RedisURL)
	cancel()
	if err != nil {
		log.Fatalf("[gateway] redis: %v", err)
	}
	defer rdb.Close()
	log.Println("[gateway] redis connected")

	// ── gRPC Clients ───────────────────────────────────────────
	clients, err := client.New(
		cfg.UserServiceAddr,
		cfg.CatalogServiceAddr,
		cfg.OrderServiceAddr,
		cfg.PaymentServiceAddr,
	)
	if err != nil {
		log.Fatalf("[gateway] grpc clients: %v", err)
	}
	defer clients.Close()

	// ── Router ─────────────────────────────────────────────────
	router := newRouter(clients, rdb)

	// ── HTTP Server ────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("[gateway] listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[gateway] serve: %v", err)
		}
	}()

	<-quit
	log.Println("[gateway] shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[gateway] shutdown: %v", err)
	}

	log.Println("[gateway] stopped")
}
