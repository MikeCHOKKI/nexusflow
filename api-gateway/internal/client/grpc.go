// ─── gRPC Clients ──────────────────────────────────────────────
package client

import (
	"fmt"
	"log"

	pb "github.com/MikeCHOKKI/nexusflow/api-gateway/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Clients holds all gRPC service connections.
type Clients struct {
	User    pb.UserServiceClient
	Catalog pb.CatalogServiceClient
	Order   pb.OrderServiceClient
	Payment pb.PaymentServiceClient

	Conns []*grpc.ClientConn
}

// New initialises all gRPC clients from the given addresses.
// Call Close() on the returned Clients to clean up connections.
func New(userAddr, catalogAddr, orderAddr, paymentAddr string) (*Clients, error) {
	clients := &Clients{}

	conn, err := dial(userAddr)
	if err != nil {
		return nil, fmt.Errorf("user service: %w", err)
	}
	clients.Conns = append(clients.Conns, conn)
	clients.User = pb.NewUserServiceClient(conn)

	conn, err = dial(catalogAddr)
	if err != nil {
		return nil, fmt.Errorf("catalog service: %w", err)
	}
	clients.Conns = append(clients.Conns, conn)
	clients.Catalog = pb.NewCatalogServiceClient(conn)

	conn, err = dial(orderAddr)
	if err != nil {
		return nil, fmt.Errorf("order service: %w", err)
	}
	clients.Conns = append(clients.Conns, conn)
	clients.Order = pb.NewOrderServiceClient(conn)

	conn, err = dial(paymentAddr)
	if err != nil {
		return nil, fmt.Errorf("payment service: %w", err)
	}
	clients.Conns = append(clients.Conns, conn)
	clients.Payment = pb.NewPaymentServiceClient(conn)

	log.Printf("[grpc] clients initialised: user=%s catalog=%s order=%s payment=%s",
		userAddr, catalogAddr, orderAddr, paymentAddr)

	return clients, nil
}

// Close gracefully closes all gRPC connections.
func (c *Clients) Close() {
	for _, conn := range c.Conns {
		if err := conn.Close(); err != nil {
			log.Printf("[grpc] close error: %v", err)
		}
	}
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
