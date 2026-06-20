package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	nexusflow "github.com/MikeCHOKKI/nexusflow/service-user/gen"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/config"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/db"
	grpcserver "github.com/MikeCHOKKI/nexusflow/service-user/internal/grpc"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/middleware"
	"github.com/MikeCHOKKI/nexusflow/service-user/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("=== User Service — NexusFlow ===")

	// ─── Configuration ──────────────────────────────────────────
	cfg := config.Load()
	log.Printf("config: grpc_port=%s", cfg.GRPCPort)

	// ─── PostgreSQL ─────────────────────────────────────────────
	pg, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer pg.Close()
	log.Println("postgres: connected")

	// ─── Redis ──────────────────────────────────────────────────
	rdb, err := db.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer rdb.Close()
	log.Println("redis: connected")

	// ─── Services ───────────────────────────────────────────────
	authSvc := service.NewAuthService(pg, rdb, cfg.JWTSecret)

	// ─── Serveur gRPC ───────────────────────────────────────────
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.UnaryPanicRecoveryInterceptor,
			middleware.UnaryLoggingInterceptor,
		),
	)

	userSrv := grpcserver.NewUserServer(authSvc, pg)
	nexusflow.RegisterUserServiceServer(grpcServer, userSrv)

	// Activer la reflection gRPC (utile pour grpcurl, Postman…).
	reflection.Register(grpcServer)

	// ─── Démarrage ──────────────────────────────────────────────
	go func() {
		log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	// ─── Graceful shutdown ──────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("signal %v received, shutting down…", sig)

	grpcServer.GracefulStop()
	log.Println("server stopped")
}
