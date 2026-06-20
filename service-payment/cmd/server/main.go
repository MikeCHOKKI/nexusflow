// ─── Payment Service — Entry Point ─────────────────────────────
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	nexusflow "github.com/MikeCHOKKI/nexusflow/protos/gen/go/nexusflow"

	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/config"
	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/db"
	internalgrpc "github.com/MikeCHOKKI/nexusflow/service-payment/internal/grpc"
	"github.com/MikeCHOKKI/nexusflow/service-payment/internal/payment"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("┌──────────────────────────────────────────────┐")
	log.Println("│  NexusFlow — Payment Service                │")
	log.Println("└──────────────────────────────────────────────┘")

	// ── Load configuration ──────────────────────────────────────
	cfg := config.Load()

	// ── Initialise Redis ────────────────────────────────────────
	redisClient, err := db.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("✓ Connected to Redis at %s", cfg.RedisURL)

	// ── Initialise payment provider (stub) ──────────────────────
	provider := payment.NewStubProvider(redisClient)
	log.Println("✓ Payment provider initialised (stub mode)")

	// ── Start gRPC server ───────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()

	// Register the PaymentService
	paymentServer := internalgrpc.NewPaymentServer(provider)
	nexusflow.RegisterPaymentServiceServer(grpcServer, paymentServer)

	// Enable reflection for debugging (grpcurl, etc.)
	reflection.Register(grpcServer)

	// ── Graceful shutdown ───────────────────────────────────────
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		log.Printf("⏳ Received signal %v, shutting down...", sig)
		grpcServer.GracefulStop()
	}()

	log.Printf("🚀 Payment service listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}
