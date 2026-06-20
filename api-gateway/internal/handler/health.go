// ─── Health Check Handler ──────────────────────────────────────
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/MikeCHOKKI/nexusflow/api-gateway/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// HealthHandler exposes the /health endpoint.
type HealthHandler struct {
	clients *client.Clients
	conns   []*grpc.ClientConn
}

// NewHealthHandler creates a handler that reports gateway + upstream status.
func NewHealthHandler(clients *client.Clients, conns []*grpc.ClientConn) *HealthHandler {
	return &HealthHandler{clients: clients, conns: conns}
}

// Check responds to GET /health with gateway and per-service status.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	_ = ctx // kept for future gRPC health checks

	services := map[string]string{}

	// Connectivity-based health: check each connection state.
	if h.clients != nil {
		services["user"] = connState(h.conns[0])
		services["catalog"] = connState(h.conns[1])
		services["order"] = connState(h.conns[2])
		services["payment"] = connState(h.conns[3])
	}

	overall := "ok"
	for _, s := range services {
		if s != "ok" {
			overall = "degraded"
			break
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":   overall,
		"services": services,
	})
}

func connState(conn *grpc.ClientConn) string {
	switch conn.GetState() {
	case connectivity.Ready:
		return "ok"
	case connectivity.Connecting:
		return "connecting"
	case connectivity.TransientFailure:
		return "unavailable"
	default:
		return "idle"
	}
}
